package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/app"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/dynamo"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/localize"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/notify"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/onelogin"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/pay"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/place"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/secrets"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/telemetry"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/templatefn"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/uid"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/mod/sumdb/dirhash"
)

var Tag string

func main() {
	ctx := context.Background()
	logger := logging.New(os.Stdout, "opg-modernising-lpa")

	var (
		appPublicURL          = env.Get("APP_PUBLIC_URL", "http://localhost:5050")
		authRedirectBaseURL   = env.Get("AUTH_REDIRECT_BASE_URL", "http://localhost:5050")
		webDir                = env.Get("WEB_DIR", "web")
		awsBaseURL            = env.Get("AWS_BASE_URL", "")
		clientID              = env.Get("CLIENT_ID", "client-id-value")
		issuer                = env.Get("ISSUER", "http://sign-in-mock:7012")
		dynamoTableLpas       = env.Get("DYNAMODB_TABLE_LPAS", "")
		notifyBaseURL         = env.Get("GOVUK_NOTIFY_BASE_URL", "")
		notifyIsProduction    = env.Get("GOVUK_NOTIFY_IS_PRODUCTION", "") == "1"
		ordnanceSurveyBaseURL = env.Get("ORDNANCE_SURVEY_BASE_URL", "http://ordnance-survey-mock:4011")
		payBaseURL            = env.Get("GOVUK_PAY_BASE_URL", "http://pay-mock:4010")
		port                  = env.Get("APP_PORT", "8080")
		yotiClientSdkID       = env.Get("YOTI_CLIENT_SDK_ID", "")
		yotiScenarioID        = env.Get("YOTI_SCENARIO_ID", "")
		yotiSandbox           = env.Get("YOTI_SANDBOX", "") == "1"
		xrayEnabled           = env.Get("XRAY_ENABLED", "") == "1"
		rumConfig             = page.RumConfig{
			GuestRoleArn:      env.Get("AWS_RUM_GUEST_ROLE_ARN", ""),
			Endpoint:          env.Get("AWS_RUM_ENDPOINT", ""),
			ApplicationRegion: env.Get("AWS_RUM_APPLICATION_REGION", ""),
			IdentityPoolID:    env.Get("AWS_RUM_IDENTITY_POOL_ID", ""),
			ApplicationID:     env.Get("AWS_RUM_APPLICATION_ID", ""),
		}
		uidBaseURL  = env.Get("UID_BASE_URL", "http://uid-mock:8080")
		metadataURL = env.Get("ECS_CONTAINER_METADATA_URI_V4", "")
	)

	staticHash, err := dirhash.HashDir(webDir+"/static", webDir, dirhash.DefaultHash)
	if err != nil {
		logger.Fatal(err)
	}
	staticHash = url.QueryEscape(staticHash[3:11])

	httpClient := &http.Client{Timeout: 10 * time.Second}

	if xrayEnabled {
		resource, err := ecs.NewResourceDetector().Detect(ctx)
		if err != nil {
			logger.Fatal(err)
		}

		shutdown, err := telemetry.Setup(ctx, resource)
		if err != nil {
			logger.Fatal(err)
		}
		defer shutdown(ctx)

		httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)
	}

	var region string
	if metadataURL != "" {
		region, err = awsRegion(metadataURL)
		if err != nil {
			logger.Print("error getting region:", err)
		}
	}

	tmpls, err := template.Parse(webDir+"/template", templatefn.All(Tag, region))
	if err != nil {
		logger.Fatal(err)
	}

	bundle := localize.NewBundle("lang/en.json", "lang/cy.json")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to load SDK config: %w", err))
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	if len(awsBaseURL) > 0 {
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsBaseURL,
				SigningRegion: "eu-west-1",
			}, nil
		})
	}

	dynamoClient, err := dynamo.NewClient(cfg, dynamoTableLpas)
	if err != nil {
		logger.Fatal(err)
	}

	secretsClient, err := secrets.NewClient(cfg, time.Hour)
	if err != nil {
		logger.Fatal(err)
	}

	sessionKeys, err := secretsClient.CookieSessionKeys(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	sessionStore := sessions.NewCookieStore(sessionKeys...)

	redirectURL := authRedirectBaseURL + page.Paths.AuthRedirect

	signInClient, err := onelogin.Discover(ctx, logger, httpClient, secretsClient, issuer, clientID, redirectURL)
	if err != nil {
		logger.Fatal(err)
	}

	payApiKey, err := secretsClient.Secret(ctx, secrets.GovUkPay)
	if err != nil {
		logger.Fatal(err)
	}

	payClient := &pay.Client{
		BaseURL:    payBaseURL,
		ApiKey:     payApiKey,
		HttpClient: httpClient,
	}

	yotiPrivateKey, err := secretsClient.SecretBytes(ctx, secrets.YotiPrivateKey)
	if err != nil {
		logger.Fatal(err)
	}

	yotiClient, err := identity.NewYotiClient(yotiScenarioID, yotiClientSdkID, yotiPrivateKey)
	if err != nil {
		logger.Fatal(err)
	}
	if yotiSandbox {
		if err := yotiClient.SetupSandbox(); err != nil {
			logger.Fatal(err)
		}
	}

	osApiKey, err := secretsClient.Secret(ctx, secrets.OrdnanceSurvey)
	if err != nil {
		logger.Fatal(err)
	}

	addressClient := place.NewClient(ordnanceSurveyBaseURL, osApiKey, httpClient)

	notifyApiKey, err := secretsClient.Secret(ctx, secrets.GovUkNotify)
	if err != nil {
		logger.Fatal(err)
	}

	notifyClient, err := notify.New(notifyIsProduction, notifyBaseURL, notifyApiKey, httpClient)
	if err != nil {
		logger.Fatal(err)
	}

	uidClient := uid.New(uidBaseURL, httpClient, cfg, v4.NewSigner(), time.Now)

	mux := http.NewServeMux()
	mux.HandleFunc(page.Paths.HealthCheck.Service, func(w http.ResponseWriter, r *http.Request) {})
	mux.Handle(page.Paths.HealthCheck.Dependency, page.DependencyHealthCheck(logger, uidClient))
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webDir+"/robots.txt")
	})
	mux.Handle("/static/", http.StripPrefix("/static", handlers.CompressHandler(page.CacheControlHeaders(http.FileServer(http.Dir(webDir+"/static/"))))))
	mux.Handle(page.Paths.AuthRedirect, page.AuthRedirect(logger, sessionStore))
	mux.Handle(page.Paths.YotiRedirect, page.YotiRedirect(logger, sessionStore))
	mux.Handle(page.Paths.CookiesConsent, page.CookieConsent(page.Paths))
	mux.Handle("/cy/", http.StripPrefix("/cy", app.App(logger, bundle.For(localize.Cy), localize.Cy, tmpls, sessionStore, dynamoClient, appPublicURL, payClient, yotiClient, notifyClient, addressClient, rumConfig, staticHash, page.Paths, signInClient, uidClient)))
	mux.Handle("/", app.App(logger, bundle.For(localize.En), localize.En, tmpls, sessionStore, dynamoClient, appPublicURL, payClient, yotiClient, notifyClient, addressClient, rumConfig, staticHash, page.Paths, signInClient, uidClient))

	var handler http.Handler = mux
	if xrayEnabled {
		handler = telemetry.WrapHandler(mux)
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           page.Recover(tmpls.Get("error-500.gohtml"), logger, bundle, handler),
		ReadHeaderTimeout: 20 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Print("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Print("signal received: ", sig)

	tc, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Print(err)
	}
}

func awsRegion(metadataURL string) (string, error) {
	resp, err := http.Get(metadataURL + "/task")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var metadata struct{ TaskARN string }
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return "", err
	}

	parts := strings.Split(metadata.TaskARN, ":")
	if len(parts) < 4 {
		return "", fmt.Errorf("TaskARN contained only %d parts", len(parts))
	}

	return parts[3], nil
}
