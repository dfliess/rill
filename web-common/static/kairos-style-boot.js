// Kairos Styles (kairos-cloud#242): paints the Style before the app loads, from
// the viewer's choice or the last default this browser saw for the Org in the
// URL, so a reload never flashes another Style. style-control.ts takes over
// once the app runs. A file rather than inline because the CSP hashes only
// the scripts SvelteKit generates. Keep STYLES in step with styles.ts
// (scripts/style-guard.js checks it).
(function () {
  var STYLES = ["kairos-v1", "kairos", "aurora"];
  try {
    var path = window.location.pathname;
    var org = path.split("/")[1];
    var user = path.indexOf("/-/embed") === -1 ? localStorage.getItem("kairos:style") : null;
    var cached = org && org !== "-" ? localStorage.getItem("kairos:style:org:" + org) : null;
    var style = [user, cached].filter(function (s) {
      return STYLES.indexOf(s) !== -1;
    })[0];
    if (style) document.documentElement.setAttribute("data-kairos-style", style);
  } catch (e) {
    // Storage blocked: kairos.css also applies without the attribute.
  }
})();
