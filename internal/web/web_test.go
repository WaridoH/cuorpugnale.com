package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomePageRendersSocialLinks(t *testing.T) {
	rec := renderHome(t)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, link := range []string{
		defaultSpotifyURL,
		"https://youtube.com/@cuorpugnale?si=GMnp_eG1ujakmclG",
		"https://www.instagram.com/cuorpugnale",
	} {
		if !strings.Contains(body, link) {
			t.Errorf("home page does not contain %q", link)
		}
	}
}

func TestHomePageRendersCampaignLedStructure(t *testing.T) {
	rec := renderHome(t)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := strings.Join(strings.Fields(rec.Body.String()), " ")
	for _, snippet := range []string{
		`<main id="contenuto">`,
		`<nav class="site-nav" aria-label="Navigazione principale">`,
		`href="#avventura"`,
		`href="#party"`,
		`href="#ascolta"`,
		`href="#chi-siamo"`,
		`<section class="campaign-hero" aria-labelledby="page-title">`,
		`<span>Dimora</span> Age of Umbra`,
		`class="placeholder-label">Sinossi in arrivo</span>`,
		`class="placeholder-label">Schede in arrivo</span>`,
		`class="party__gallery"`,
	} {
		if !strings.Contains(body, snippet) {
			t.Errorf("home page does not contain %q", snippet)
		}
	}
	for _, blocked := range []string{
		`class="credits"`,
		`Ogni storia ha molte mani`,
		`story-panel__number`,
		`party__grid`,
		`Personaggio 04`,
		`aria-hidden="true">I</div>`,
		`aria-hidden="true">II</div>`,
		`aria-hidden="true">III</div>`,
		`aria-hidden="true">IV</div>`,
	} {
		if strings.Contains(body, blocked) {
			t.Errorf("home page contains removed content %q", blocked)
		}
	}
	if got := strings.Count(body, `<article class="character-card">`); got != 3 {
		t.Errorf("character card count = %d, want 3", got)
	}
}

func TestHomePageRendersTeamInProjectOrder(t *testing.T) {
	rec := renderHome(t)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	previous := -1
	for _, member := range []struct {
		name  string
		image string
	}{
		{name: "Walid Husein", image: "walid-husein.jpg"},
		{name: "Gianluca Danaro", image: "gianluca-danaro.jpg"},
		{name: "Lorenzo Magalotti", image: "lorenzo-magalotti.jpg"},
		{name: "Supernova Collective", image: "supernova-collective.jpg"},
		{name: "Lorenzo Morelli", image: "lorenzo-morelli.jpg"},
		{name: "Federica Passi", image: "federica-passi.jpg"},
		{name: "Denise Venanzetti", image: "denise-venanzetti.jpg"},
	} {
		snippet := `src="/static/img/team/` + member.image + `"`
		position := strings.Index(body, snippet)
		if position == -1 {
			t.Errorf("home page does not contain team member %q", member.name)
			continue
		}
		if position < previous {
			t.Errorf("team member %q appears out of order", member.name)
		}
		previous = position
	}
}

func TestHomePageRendersSupernovaCollective(t *testing.T) {
	rec := renderHome(t)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := strings.Join(strings.Fields(rec.Body.String()), " ")
	for _, snippet := range []string{
		`src="/static/img/team/supernova-collective.jpg"`,
		`Post-produzione audio, composizione colonna sonora originale, sound design`,
		`Giovanni Mennuni, Matteo Mezzabotta e Gianluca Danaro`,
		`Hanno collaborato con: Radio Deejay (Dungeons&amp;Deejay), RaiPlay, Università RomaTre, InnTale, Sabaku No Maiku/Camilla d’Onofrio, Roberto Recchioni, Mog’s Chronicles.`,
		`https://www.instagram.com/supernova_collective/`,
		`https://linktr.ee/SupernovaCollective`,
		`Ovunque amiate ascoltare i vostri podcast`,
	} {
		if !strings.Contains(body, snippet) {
			t.Errorf("home page does not contain %q", snippet)
		}
	}
	if strings.Contains(body, `Presto ovunque amiate ascoltare i vostri podcast`) {
		t.Error("home page still describes podcast availability as forthcoming")
	}
	if got := strings.Count(body, `https://www.instagram.com/supernova_collective/`); got != 2 {
		t.Errorf("Supernova Instagram link count = %d, want 2", got)
	}
}

func TestHomePageCanOverrideSpotifyLink(t *testing.T) {
	t.Setenv("SPOTIFY_URL", "https://open.spotify.com/show/example")

	if got := spotifyURL(); got != "https://open.spotify.com/show/example" {
		t.Errorf("spotifyURL() = %q, want overridden Spotify URL", got)
	}
}

func TestHomePageUsesResponsiveCampaignArtwork(t *testing.T) {
	rec := renderHome(t)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, snippet := range []string{
		`rel="preload"`,
		`as="image"`,
		`href="/static/img/campaign/age-of-umbra-1920.avif"`,
		`imagesizes="100vw"`,
		`<picture class="campaign-hero__artwork">`,
		`type="image/avif"`,
		`/static/img/campaign/age-of-umbra-960.avif 960w`,
		`/static/img/campaign/age-of-umbra-2560.avif 2560w`,
		`src="/static/img/campaign/age-of-umbra-1920.jpg"`,
		`alt="Il party di Dimora: Age of Umbra davanti alla città nella nebbia"`,
		`width="1920"`,
		`height="1080"`,
		`property="og:image" content="https://cuorpugnale.com/static/img/campaign/age-of-umbra-1920.jpg"`,
		`<link rel="stylesheet" href="/static/css/age-of-umbra-v1.css" />`,
	} {
		if !strings.Contains(body, snippet) {
			t.Errorf("home page does not contain %q", snippet)
		}
	}
}

func TestHomePagePrioritizesListening(t *testing.T) {
	rec := renderHome(t)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, snippet := range []string{
		`<section class="listen" id="ascolta" aria-labelledby="listen-title">`,
		`href="#ascolta"`,
		`Ascolta ora`,
		`data-testid="embed-iframe"`,
		`class="listen__spotify"`,
		`title="Dimora: Age of Umbra su Spotify"`,
		`src="https://open.spotify.com/embed/show/033nh6SpB5aFTDuQ6FMkSI?utm_source=generator"`,
		`allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"`,
	} {
		if !strings.Contains(body, snippet) {
			t.Errorf("home page does not contain %q", snippet)
		}
	}
}

func TestCampaignAssetsAreServed(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	handler := server.Handler()
	for _, path := range []string{
		"/static/css/age-of-umbra-v1.css",
		"/static/img/campaign/age-of-umbra-960.avif",
		"/static/img/campaign/age-of-umbra-1920.avif",
		"/static/img/campaign/age-of-umbra-2560.avif",
		"/static/img/campaign/age-of-umbra-1920.jpg",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("%s status = %d, want %d", path, rec.Code, http.StatusOK)
		}
		if cacheControl := rec.Header().Get("Cache-Control"); !strings.Contains(cacheControl, "immutable") {
			t.Errorf("%s Cache-Control = %q, want immutable static cache", path, cacheControl)
		}
	}
}

func TestCampaignStylesIncludeGalleryAndColorPortraitInteraction(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/static/css/age-of-umbra-v1.css", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	for _, snippet := range []string{
		`.party__gallery {`,
		`overflow-x: auto;`,
		`scroll-snap-type: x mandatory;`,
		`.person:hover .person__photo,`,
		`.person:focus-within .person__photo {`,
		`filter: grayscale(0);`,
	} {
		if !strings.Contains(rec.Body.String(), snippet) {
			t.Errorf("campaign stylesheet does not contain %q", snippet)
		}
	}
}

func TestSupernovaImageIsServed(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	path := "/static/img/team/supernova-collective.jpg"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("%s status = %d, want %d", path, rec.Code, http.StatusOK)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "image/jpeg" {
		t.Errorf("%s Content-Type = %q, want image/jpeg", path, contentType)
	}
	if cacheControl := rec.Header().Get("Cache-Control"); !strings.Contains(cacheControl, "immutable") {
		t.Errorf("%s Cache-Control = %q, want immutable static cache", path, cacheControl)
	}
}

func TestHomePageDoesNotLoadExternalFonts(t *testing.T) {
	rec := renderHome(t)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, blocked := range []string{
		"fonts.googleapis.com",
		"fonts.gstatic.com",
		"Playfair Display",
		"Lora",
	} {
		if strings.Contains(body, blocked) {
			t.Errorf("home page contains external font dependency %q", blocked)
		}
	}
}

func TestSecurityPolicyDoesNotAllowExternalFonts(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	csp := rec.Header().Get("Content-Security-Policy")
	for _, blocked := range []string{
		"fonts.googleapis.com",
		"fonts.gstatic.com",
	} {
		if strings.Contains(csp, blocked) {
			t.Errorf("Content-Security-Policy contains external font source %q", blocked)
		}
	}
}

func TestSecurityPolicyDoesNotAllowInlineScripts(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if csp := rec.Header().Get("Content-Security-Policy"); strings.Contains(csp, "'unsafe-inline'") {
		t.Errorf("Content-Security-Policy still allows inline scripts: %q", csp)
	}
}

func renderHome(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()

	server, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	return rec
}
