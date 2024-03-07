package main

import (
	"encoding/json"
	"fmt"

	ep "github.com/ogame-ninja/extension-patcher"
)

func main() {
	const (
		extensionName         = "infinity"
		webstoreURL           = "https://chrome.google.com/webstore/detail/ogame-infinity/hfojakphgokgpbnejoobfamojbgolcbo"
		infinity_2_4_1_sha256 = "87cb6e2d49d5a31c263e24314ff9a4eb7bcb3111095c0392571ed7cf473ff35f"
	)

	files := []ep.FileAndProcessors{
		ep.NewFile("/manifest.json", processManifest),
		ep.NewFile("/ctxcontent/index.js", processCtXContextIndexJS),
		ep.NewFile("/ctxcontent/data-helper.js", processDataHelperJS),
		ep.NewFile("/ctxcontent/helpers/universe.alliances.js", processHelpersUniverseAlliancesJS),
		ep.NewFile("/ctxcontent/helpers/universe.players.js", processHelpersUniversePlayersJS),
		ep.NewFile("/ctxcontent/helpers/universe.planets.js", processHelpersUniversePlanetsJS),
		ep.NewFile("/ctxcontent/helpers/universe.highscore.js", processHelpersUniverseHighscoreJS),
		ep.NewFile("/ctxcontent/services/request.ogameAlliances.js", processServiceOgameAlliancesJS),
		ep.NewFile("/ctxcontent/services/request.ogamePlayers.js", processServiceOgamePlayersJS),
		ep.NewFile("/ctxcontent/services/request.ogamePlanets.js", processServiceOgamePlanetsJS),
		ep.NewFile("/ctxcontent/services/request.ogameHighscore.js", processServiceOgameHighscoreJS),
		ep.NewFile("/ogkush.js", processOgkushJS),
		ep.NewFile("/background.js", processBackgroundJS),
	}

	p, err := ep.New(ep.Params{
		ExtensionName:  extensionName,
		ExpectedSha256: infinity_2_4_1_sha256,
		WebstoreURL:    webstoreURL,
		Files:          files,
	})
	if err != nil {
		panic(err)
	}
	p.Start()
}

var replN = ep.MustReplaceN

type Manifest struct {
	Action struct {
		DefaultIcon struct {
			Num32 string `json:"32"`
		} `json:"default_icon"`
		DefaultTitle string `json:"default_title"`
	} `json:"action"`
	Background struct {
		ServiceWorker string `json:"service_worker"`
	} `json:"background"`
	ContentScripts []struct {
		CSS            []string `json:"css"`
		ExcludeMatches []string `json:"exclude_matches"`
		Js             []string `json:"js"`
		Matches        []string `json:"matches"`
		RunAt          string   `json:"run_at"`
	} `json:"content_scripts"`
	Description     string   `json:"description"`
	HostPermissions []string `json:"host_permissions"`
	Icons           struct {
		Num32  string `json:"32"`
		Num128 string `json:"128"`
		Num256 string `json:"256"`
		Num512 string `json:"512"`
	} `json:"icons"`
	ManifestVersion        int      `json:"manifest_version"`
	Name                   string   `json:"name"`
	Permissions            []string `json:"permissions"`
	UpdateURL              string   `json:"update_url"`
	Version                string   `json:"version"`
	WebAccessibleResources []struct {
		ExtensionIds []string `json:"extension_ids"`
		Matches      []string `json:"matches"`
		Resources    []string `json:"resources"`
	} `json:"web_accessible_resources"`
}

func processManifest(by []byte) []byte {
	var data Manifest
	if err := json.Unmarshal(by, &data); err != nil {
		fmt.Println(err)
	}
	data.Name = "Ogame Infinity Ninja"
	data.ContentScripts[0].Matches = append(data.ContentScripts[0].Matches, "*://*/bots/*/browser/html/*")
	data.HostPermissions = append(data.HostPermissions, "<all_urls>")
	out, _ := json.MarshalIndent(data, "", "  ")
	return out
}

func processCtXContextIndexJS(by []byte) []byte {
	by = replN(by, `const UNIVERSE=window.location.host.split(".")[0];`,
		`const universeNum = /browser\/html\/s(\d+)-(\w+)/.exec(window.location.href)[1];
const lang = /browser\/html\/s(\d+)-(\w+)/.exec(window.location.href)[2];
const UNIVERSE = "s" + universeNum + "-" + lang;
const PROTOCOL = window.location.protocol;
const HOST = window.location.host;`, 1)
	by = replN(by, `universe:UNIVERSE`, `protocol:window.location.protocol,host:window.location.host,universe: UNIVERSE`, 1)
	by = replN(by, `new DataHelper(UNIVERSE)`, `new DataHelper(PROTOCOL, HOST, UNIVERSE)`, 2)
	return by
}

func processDataHelperJS(by []byte) []byte {
	by = replN(by, `constructor(universe){`,
		`constructor(protocol,host,universe){
this.protocol=protocol;
this.host=host;
this.universeNum=/s(\d+)-(\w+)/.exec(universe)[1];
this.universeLang=/s(\d+)-(\w+)/.exec(universe)[2];`, 1)
	by = replN(by, `getPlayersHighscore(this.universe)`, `getPlayersHighscore(this.protocol, this.host, this.universe, this.universeNum, this.universeLang)`, 1)
	by = replN(by, `getPlayers(this.universe)`, `getPlayers(this.protocol, this.host, this.universe, this.universeNum, this.universeLang)`, 1)
	by = replN(by, `getPlanets(this.universe)`, `getPlanets(this.protocol, this.host, this.universe, this.universeNum, this.universeLang)`, 1)
	by = replN(by, `getAlliances(this.universe)`, `getAlliances(this.protocol, this.host, this.universe, this.universeNum, this.universeLang)`, 1)
	return by
}

func processHelpersUniverseAlliancesJS(by []byte) []byte {
	by = replN(by, `function getAlliances(universe)`, `function getAlliances(protocol, host, universe, universeNum, universeLang)`, 1)
	by = replN(by, `requestOGameAlliances(universe)`, `requestOGameAlliances(protocol, host, universe, universeNum, universeLang)`, 1)
	return by
}

func processHelpersUniversePlayersJS(by []byte) []byte {
	by = replN(by, `function getPlayers(universe)`, `function getPlayers(protocol, host, universe, universeNum, universeLang)`, 1)
	by = replN(by, `requestOGamePlayers(universe)`, `requestOGamePlayers(protocol, host, universe, universeNum, universeLang)`, 1)
	return by
}

func processHelpersUniversePlanetsJS(by []byte) []byte {
	by = replN(by, `function getPlanets(universe)`, `function getPlanets(protocol, host, universe, universeNum, universeLang)`, 1)
	by = replN(by, `requestOGamePlanets(universe)`, `requestOGamePlanets(protocol, host, universe, universeNum, universeLang)`, 1)
	return by
}

func processHelpersUniverseHighscoreJS(by []byte) []byte {
	by = replN(by, `function requestHighscore(universe,category)`, `function requestHighscore(protocol, host, universe, category, universeNum, universeLang)`, 1)
	by = replN(by, `getPlayersHighscore(universe)`, `getPlayersHighscore(protocol, host, universe, universeNum, universeLang)`, 1)
	by = replN(by, `requestHighscore(universe,HIGHSCORE_CATEGORY.PLAYER)`, `requestHighscore(protocol, host, universe, HIGHSCORE_CATEGORY.PLAYER, universeNum, universeLang)`, 1)
	by = replN(by, `requestHighscore(universe,HIGHSCORE_CATEGORY.ALLIANCE)`, `requestHighscore(protocol, host, universe, HIGHSCORE_CATEGORY.ALLIANCE, universeNum, universeLang)`, 1)
	by = replN(by, `requestOGameHighScore(universe,category,type)`, `requestOGameHighScore(protocol, host, universe, category, type, universeNum, universeLang)`, 1)
	return by
}

func processServiceOgameAlliancesJS(by []byte) []byte {
	by = replN(by, `function requestOGameAlliances(universe)`, `function requestOGameAlliances(protocol, host, universe, universeNum, universeLang)`, 1)
	by = replN(by, `https://${universe}.ogame.gameforge.com/api/alliances.xml`,
		`${protocol}//${host}/api/s${universeNum}/${universeLang}/alliances.xml`, 1)
	return by
}

func processServiceOgamePlayersJS(by []byte) []byte {
	by = replN(by, `function requestOGamePlayers(universe)`, `function requestOGamePlayers(protocol, host, universe, universeNum, universeLang)`, 1)
	by = replN(by, `https://${universe}.ogame.gameforge.com/api/players.xml`,
		`${protocol}//${host}/api/s${universeNum}/${universeLang}/players.xml`, 1)
	return by
}

func processServiceOgamePlanetsJS(by []byte) []byte {
	by = replN(by, `function requestOGamePlanets(universe)`, `function requestOGamePlanets(protocol, host, universe, universeNum, universeLang)`, 1)
	by = replN(by, `https://${universe}.ogame.gameforge.com/api/universe.xml`,
		`${protocol}//${host}/api/s${universeNum}/${universeLang}/universe.xml`, 1)
	return by
}

func processServiceOgameHighscoreJS(by []byte) []byte {
	by = replN(by, `function requestOGameHighScore(universe,category,type)`, `function requestOGameHighScore(protocol, host, universe, category, type, universeNum, universeLang)`, 1)
	by = replN(by, `https://${universe}.ogame.gameforge.com/api/highscore.xml`,
		`${protocol}//${host}/api/s${universeNum}/${universeLang}/highscore.xml`, 1)
	return by
}

func processOgkushJS(by []byte) []byte {
	by = append([]byte(`const universeNum = /browser\/html\/s(\d+)-(\w+)/.exec(window.location.href)[1];
const lang = /browser\/html\/s(\d+)-(\w+)/.exec(window.location.href)[2];
const UNIVERSE = "s" + universeNum + "-" + lang;
const PLAYER_ID = document.querySelector("meta[name=ogame-player-id]").content;
const localStoragePrefix = UNIVERSE + "-" + PLAYER_ID + "-";`), by...)
	by = replN(by, `localStorage.getItem(`, `localStorage.getItem(localStoragePrefix+`, 2)
	by = replN(by, `localStorage.setItem(`, `localStorage.setItem(localStoragePrefix+`, 6)
	by = replN(by, `window.location.host.replace(/\D/g,"");`, `universeNum;`, 1)
	by = replN(by, `https://s${this.universe}-${this.gameLang}.ogame.gameforge.com/api/serverData.xml`, `/api/s${universeNum}/${lang}/serverData.xml`, 1)
	by = replN(by, `https://s${this.universe}-${this.gameLang}.ogame.gameforge.com/game/index.php`, ``, 16)
	by = replN(by, `;for(var x in localStorage){`, `;for(var x in localStorage){if(!x.startsWith(localStoragePrefix)){continue;}`, 1)
	by = replN(by, `purgeLocalStorage(){for(var x in localStorage){if(x!="ogk-data"){`, `purgeLocalStorage(){for(var x in localStorage){if(!x.startsWith(localStoragePrefix)){continue;}if(x!=localStoragePrefix+"ogk-data"){`, 1)
	by = replN(by, `document.location.origin+"/game/index.php`, `"`, 2)
	by = replN(by, "`/game/index.php?", "`?", 2)
	by = replN(by, `"/game/index.php?`, `"?`, 1)
	by = replN(by, `"https://"+window.location.host+window.location.pathname+`, ``, 14)
	return by
}

func processBackgroundJS(by []byte) []byte {
	return by
}
