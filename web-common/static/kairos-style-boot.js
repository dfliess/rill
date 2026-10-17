// Kairos Styles (kairos-cloud#242): paints the Style before the app loads, from
// the embed parameter, the viewer's choice or the last default this browser
// saw for the Org in the URL, so a reload never flashes another Style. style-control.ts takes over
// once the app runs. A file rather than inline because the CSP hashes only
// the scripts SvelteKit generates. Keep STYLES in step with styles.ts
// (scripts/style-guard.js checks it).
(function () {
  var STYLES = ["kairos", "aurora", "nova", "editorial", "lumen", "glass", "apetic", "base"];
  try {
    var path = window.location.pathname;
    var org = path.split("/")[1];
    // An embed paints its URL parameter (GetIFrame puts the Org's there), never the viewer's choice.
    var embed = path.indexOf("/-/embed") !== -1;
    var param = embed ? new URLSearchParams(window.location.search).get("kairos_style") : null;
    var user = embed ? null : localStorage.getItem("kairos:style");
    var cached = org && org !== "-" ? localStorage.getItem("kairos:style:org:" + org) : null;
    var style = [param, user, cached].filter(function (s) {
      return STYLES.indexOf(s) !== -1;
    })[0];
    if (style) document.documentElement.setAttribute("data-kairos-style", style);
  } catch (e) {
    // Storage blocked: kairos.css also applies without the attribute.
  }
})();
