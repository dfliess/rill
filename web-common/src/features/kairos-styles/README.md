# Kairos Styles (Appearances)

An Appearance is the look of the whole app (surfaces, text, borders, shape,
typography), in a light and a dark version. It sits between Rill's Mode
(light/dark/system) and the project theme (`themes/*.yaml` in a project).
Design and decisions: kairos-cloud#242 and ADR-0035.

## Adding an Appearance

1. A CSS file here, scoped to `:root[data-kairos-style="<id>"]`, with explicit
   `&:not(.dark)` and `&.dark` blocks (they must beat upstream `:root` and
   `:root.dark` on specificity, not on load order).
2. Its id in `STYLES` (`styles.ts`), in `kairosStyles`
   (`admin/server/kairos_styles.go`, what an Org may pick) and in
   `web-common/static/kairos-style-boot.js`, plus its label in the two
   pickers (`StyleMenu.svelte`, `OrgStyleSettings.svelte`) and the i18n
   messages.
   `scripts/style-guard.js` fails if the three lists drift.
3. Its expected computed values in `scripts/style-smoke.js`.
4. Fonts go in `web-common/static/fonts`: the CSP only allows `font-src 'self'`,
   so Google Fonts never load (a mockup that links them renders fallbacks).
5. A header whose background differs from the mode (a dark band in light mode,
   as in `kairos`, `nova` and `apetic`) shows the logos drawn for it, inside
   the same rule that paints it: hide `kairos-logo-on-light` and
   `kairos-org-logo-on-light`, show `kairos-org-logo-on-dark`. Guard the rule
   with `:not(:has(.kairos-org-logo))`: an Org with only a light logo gets no
   band. A Style that does nothing here gets the mode's default, set in
   `HeaderLogo.svelte`.

   Why there and not a declared header tone (say, a `headerTone` in
   `styles.ts` that `HeaderLogo` reads): the band rule is already the one place
   a Style says its header is dark, so the logo switch can't drift from it; a
   tone in TS would be a second source of truth beside the band's CSS.
   `scripts/style-smoke.js` checks the shown logo (`orgLogo`) per Style and mode.

## How it mixes with the project theme

The project theme is Rill's, and it only applies inside the dashboard:
`Theme.generateCSS` (`features/themes/theme.ts`) writes its tokens on
`.dashboard-theme-boundary`, so the header, menus, chat and the Kairos Home
boxes never see it. Inside the boundary it re-declares whatever the theme sets,
and that wins over the Appearance:

- From just `primary`/`secondary` it builds `--color-primary-*`,
  `--color-secondary-*`, `--color-theme-*`, and derives the semantic tokens in
  `PALETTE_DERIVED_TOKENS`: `fg-accent`, `accent-primary(-action)`,
  `accent-secondary(-action)`, `icon-accent`, `ring-focus`, `surface-hover`,
  `surface-active`, `popover-accent`, and the dimension/measure chip tokens
  (upstream #9811).
- `variables:` in the theme can override any token, surfaces and borders
  included.
- Chart palettes and `kpi-positive`/`kpi-negative` come from the theme.

Consequences when writing an Appearance:

- A decoration styled with one of the derived tokens takes the project's brand
  color inside the dashboard and the Appearance's color outside it. Seen with
  GKM (primary red `#e52822`): `kairos` card titles used `--fg-accent` and
  rendered dark red on the canvas and blue on the Home boxes (now they use
  `--kairos-label`); `aurora` list markers use `--accent-primary`, so the
  agent note bullets render red.
  Use a token the theme does not touch when the element must look the same
  everywhere; use a derived one only when it should follow the project brand.
  `scripts/style-smoke.js` checks this for `kairos` with a themed fixture.
- Tokens the theme does not set (surfaces, `border`, `fg-primary`, ...) stay
  the Appearance's, unless the project sets them through `variables:`.
- Rill components also use raw Tailwind scales the Appearance may not cover.
  The time chip in canvas filters uses `border-gray-200`: on a dark card it
  can vanish (contrast 1.04 in the retired `kairos-v1`, 1.35 in `aurora`).
  `kairos` draws it with `--border`. Check filters and chips in both modes.
