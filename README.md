# RoleWatch CLI

A small Go CLI over the [rolewatch](https://rolewatch.irfanmiral.com) job API — built so an
**agent** (or a human) can search, open, and apply to jobs from the terminal,
without a browser. Authenticate once with a personal API key; the key is stored
in `~/.rolewatch/creds.json`.

> A fork of [strelov1/freehire-cli](https://github.com/strelov1/freehire-cli) (MIT), adapted for the self-hosted RoleWatch platform. Original copyright is kept in `LICENSE`.

## Install

**curl** (prebuilt binary, no Go needed):

```bash
curl -fsSL https://rolewatch.irfanmiral.com/install.sh | sh
```

**Go:**

```bash
curl -fsSL https://rolewatch.irfanmiral.com/install.sh | sh   # installs the `rolewatch` binary
```

## Authenticate

```bash
rolewatch auth login --oauth                # opens a browser to sign in — no key to paste or manage
rolewatch auth login --token fhk_xxxxxxxx   # or: an API key from the web app (account menu → API keys)
rolewatch auth status                       # Authenticated as you@example.com @ https://rolewatch.irfanmiral.com
rolewatch auth logout                       # removes ~/.rolewatch/creds.json
```

`--oauth` registers this CLI with rolewatch's OAuth 2.1 authorization server
(rolewatch#3114), opens your browser to approve the connection, and stores the
resulting token exactly like a pasted key — nothing else in this CLI, or in
[rolewatch-mcp](https://github.com/irfanMiral/rolewatch-mcp) reading the same
`~/.rolewatch/creds.json`, needs to know which method produced it. Revoke a
browser-granted connection any time from the web app's **Connected devices**
list, without touching an API key.

`auth login` validates the credential against the API before saving, so a bad
key or a failed sign-in is never stored. Omit `--token` to be prompted on
stdin.

## Use

```bash
rolewatch profile                                               # your saved roles, skills, geography and CV (no contacts)
rolewatch experience list                                       # your whole experience bank (employments + achievements)
rolewatch experience employments add --company Acme --role SWE  # record a new job or project
rolewatch experience atoms add --claim "Cut latency 20s to 1s" --metric "20s->1s" --skill go  # record a new achievement
rolewatch experience atoms update <id> --claim "Cut p99 20s to 900ms"  # fix one field; the rest is kept
rolewatch experience employments update <id> --end "Dec 2025"   # same, for a place
rolewatch experience atoms rm <id>                              # remove an achievement (asks first)
rolewatch experience employments rm <id>                        # remove a place — it must be empty
rolewatch cv tailor <slug>                                      # start (or reopen) tailoring for a vacancy — prints the CV id
rolewatch cv list                                               # the tailored CVs you already have, with their vacancies
rolewatch cv context <cv-id>                                    # the fit analysis to reframe toward (missing_have vs missing_gap)
rolewatch cv edit <cv-id> --set 'summary=…'                     # edit by path; a claim needs --evidence <atom-id>
rolewatch cv render <cv-id> --out cv.pdf                        # download the ATS PDF
rolewatch auto-apply status <job-slug>                          # where an attempt already queued on the site stands
rolewatch auto-apply tailor <queue-id>                          # (re)start writing the CV for a queued attempt
rolewatch auto-apply review <queue-id> --approve                # queue it for unattended submission
rolewatch auto-apply review <queue-id> --decline                # park it for good
rolewatch facets                                                # list every filter's live values + counts (what to filter by)
rolewatch search "golang"                                       # list matching jobs (title · company · location · slug)
rolewatch search "backend" --remote --region eu --company acme  # facet filters (repeatable: --region, --company)
rolewatch search "data" --country BR --employment-type full_time --facet source=greenhouse  # any facet via --facet key=value
rolewatch market-fit --skills go,docker,react --category backend  # how much of the market your skills cover (+ gaps)
rolewatch market-fit --skills go --country BR                    # one skill probes that skill's demand under the filter
rolewatch job <slug>                                            # show a job's full content (incl. posting URL + company slug)
rolewatch company <slug>                                        # show a company and its open jobs
rolewatch apply <slug>                                          # mark a job applied for your account
rolewatch apply <slug> --on 2026-07-27                          # ...on the day you actually sent it
rolewatch save <slug>                                           # bookmark a job for later
rolewatch unsave <slug>                                         # remove a bookmark
rolewatch stage <slug> <stage>                                  # set application stage (applied→…→offer, or rejected/withdrawn)
rolewatch note <slug> a quick reminder                          # attach a free-text note (trailing args; no quotes needed)
rolewatch my --filter applied                                   # tracked jobs, showing stage + note (all|viewed|saved|applied)
rolewatch contribute <url>                                      # hand rolewatch a job link (see below)
rolewatch contributions                                         # boards you've contributed
rolewatch ghost report <slug> --applied-on 2026-06-01           # you applied that day and were never answered
rolewatch ghost retract <slug>                                  # withdraw that report
rolewatch inbox push < mail.json                                # upload mail your own client fetched
rolewatch inbox list --unclassified --body                      # mail still awaiting a verdict, with its text
rolewatch inbox triage <id> rejection --slug <slug>             # record what a message is + link it
rolewatch inbox list --link suggested                          # matcher proposals awaiting your word
rolewatch inbox confirm <id>                                   # accept one of them
rolewatch inbox application <id> <slug>                        # mail about an untracked application
```

**Experience bank.** `experience list` shows every achievement grouped under the
employment it belongs to (achievements with no place come back under `unplaced`) —
copy an id from there to attach a new achievement to it with `atoms add
--employment <id>`. Anything added this way is recorded with `manual` provenance,
which is what lets it later be cited as `cv edit --evidence <id>`.

`atoms update <id>` and `employments update <id>` correct an entry. Only the flags
you pass change — the API replaces the whole row, so the command reads what is
banked and sends it back with your changes applied, which is what keeps a typo fix
from deleting the metrics. `--metric` and `--skill` replace the whole list when
given.

A correction made with a key does **not** restamp provenance: an achievement the
agent inferred stays `agent_inferred`, and so stays uncitable on a CV. Only your
own edit on the site turns it into something you assert.

`atoms rm <id>` and `employments rm <id>` remove an entry. This is final — the
bank has no undo — so both ask you to type `delete` unless you pass `--yes`. A
place must be **empty** before it can go: removing one would delete every
achievement under it, and the API refuses that from a key. To retire a duplicate
place, move its achievements to the one you are keeping
(`atoms update <id> --employment <keep-id>`), then remove the shell. Folding two
achievements into one — keeping the numbers from both — is on the site.

**Tailoring a CV.** `cv tailor <job-slug>` copies your base CV into a
vacancy-bound one and prints the CV id every other `cv` command takes; it is
idempotent per vacancy, so running it again reopens the same copy. `cv list`
shows the copies you already have. It spends AI credits the first time it
creates one, and needs a résumé on the site to seed a base CV from. Editing with
a key edits as the tailoring agent: the contact block is refused, and anything
stating what you DID needs `--evidence <atom-id>` from the experience bank. The
readiness score, the revision history and undo stay on the site — deliberately,
so the agent cannot read or revert the measure of its own work.

**Auto-apply.** Queuing a *new* attempt is cookie-only on rolewatch.irfanmiral.com itself — a
real application can go out, and the browser is the only place you can watch it
happen and undo it, so there is no `auto-apply run` here and there will not be
one. Once an attempt is queued from the site, `auto-apply status <job-slug>`
shows where it stands (`tailoring` → `pending_review` → `approved`, or
`blocked`/`declined`/`failed`) and gives you the queue id the other two take.
`tailor <queue-id>` (re)starts writing the CV for it; when it reaches
`pending_review`, `status` shows exactly what would be submitted, and `review
<queue-id> --approve`/`--decline` records your decision — once only, the server
refuses a second call.

**Discovering values.** `rolewatch facets [filters]` lists every filter's live
values with a vacancy count each (and the `skills` vocabulary), so you pass real
values to `search`/`market-fit` instead of guessing. Narrow it with any filter flag
to see the values for that slice (e.g. `rolewatch facets --category backend`).

**Filters.** `search` and `market-fit` share the same market-filter flags:
`--remote`, `--region`, `--country`, `--city`, `--company`, `--category`, `--role`,
`--seniority`, `--employment-type`, `--english-level`, `--salary-min`, `--visa`
(each named facet is repeatable). Any other facet in the API's vocabulary is
reachable with `--facet key=value` (repeatable), e.g. `--facet source=greenhouse`
or `--facet company_size=startup`.

**market-fit** measures how much of the live open-vacancy market your skills reach
for a filtered role: the headline coverage (`N%` of vacancies list ≥1 of your
skills), the must-have skills you hold, and the missing skills that unlock the most
new vacancies. Pass your skills with `--skills` (comma-separated or repeated) — one
skill probes that single skill's demand; a list scores your whole stack. Note that
here `--skills` is the *measured set*, not a filter (in `search`, `--skills` filters
to jobs listing the skill).

Moderators can author postings (requires the `moderator` role; a regular key gets 403):

```bash
rolewatch jobs add --url https://acme.example/jobs/1 --title "Senior Go Developer" --company Acme --remote
rolewatch jobs add --url <url> --source workatastartup --title <t> --company <c> --description "<p>…HTML…</p>"
rolewatch jobs edit <slug> --title "Staff Go Developer"         # partial: only the flags you pass change
```

`--source` records the posting's real origin (defaults to `manual`); it does not change that
the job is flagged as manually added (that comes from the moderator authorship). `--description`
is stored and rendered as HTML, so pass HTML markup. The URL is the dedup key — re-running `add`
with the same URL updates the posting.

## Contributing a link

`rolewatch contribute <url>` hands over any vacancy or ATS board link. One sequence answers,
whichever surface you use:

| Outcome | Meaning |
|---|---|
| already in the catalogue | we carry that posting; you get its URL |
| added, company already crawled | we read the page and imported it; the rest of that company's roles follow on the next crawl |
| added, company new to us | imported, and its board is queued for crawling — **+1 AI credit** |
| queued | nothing could read the page, so a maintainer will look; not credited |

The board, not the vacancy, is what earns the credit, and only the first person to name a
board earns it — later links to the same company are still recorded (they tell us the board
matters) but pay nothing.

```bash
rolewatch contribute https://acme.recruitee.com/o/senior-go
rolewatch --json contribute https://acme.recruitee.com/o/senior-go | jq '.status, .public_slug'
rolewatch contributions
```

## A posting that may not be real

`rolewatch job <slug>` shows a signal above the description when a posting looks like it is
not being worked — how it behaves, whether the employer's own careers board carries it, and
what happened to people who applied:

```
Possibly inactive — 2 of 4 checks fired
  · Posting behaves as evergreen
  · Not on the company's own careers board (checked 2026-07-30)
  Not observed: applications here, reports from people.
```

The wording is hedged on purpose, and the unfired checks are printed too: they are what says
how sure the signal is. These are observations about a posting, never a claim about anyone's
intent.

If you applied and were never answered, say so — human evidence is the only thing that can
raise the signal past what a posting's shape can show:

```bash
rolewatch ghost report senior-go-acme --applied-on 2026-06-01
rolewatch ghost retract senior-go-acme
```

The date is the day *you* applied, and it is required. Reporting needs a verified email
address, works once per posting, and reaches no moderator — nothing here closes a job.

## Application mail — bring your own client

rolewatch does not fetch your mail. Your own client does — [himalaya], `mbsync`,
`notmuch`, anything that can emit JSON — and `inbox push` hands the result over.
Your agent then reads it and records what each message is, so the tracker stays
current without rolewatch running a mail connector or a classifier for you.

[himalaya]: https://github.com/pimalaya/himalaya

```bash
himalaya envelope list --output json | jq '[.[] | {
    external_id: .message_id, from_addr: .from.addr, from_name: .from.name,
    subject: .subject, received_at: .date }]' | rolewatch inbox push

rolewatch --json inbox list --unclassified --body --limit 50   # what still needs judging, with its text
rolewatch inbox triage 812 rejection --slug go-dev-acme-t35nijto
```

`external_id` is the deduplication key — use the message's `Message-ID`. Re-pushing
the same id updates that message instead of storing a copy, and never un-reads it,
resurrects a deleted one, or overwrites a verdict. Batches hold up to 100 messages.

`inbox list --body` returns each message's readable text (HTML-only ATS mail arrives
stripped) and marks nothing read — unlike `inbox read`, which is for opening a single
message a person asked to see. Triage signals: `acknowledgement`, `screening`,
`interview_invitation`, `assessment`, `offer`, `rejection`, `info_request`,
`incomplete_application`, `other`. A forward signal on a linked message advances that
application's stage; a settled application is never dragged back into the pipeline.

Only a deterministic match links mail automatically; anything the classifier merely
believes lands as a *suggestion* awaiting your word. `inbox list --link suggested` is
that queue, drained with `inbox confirm <id>` or `inbox reject <id>`. An unconfirmed
suggestion is a link that never happens, so a mailbox nobody confirms looks far less
matched than it is.

`inbox list --link unlinked` is the other queue: mail with no application to attach
to, usually because the application was never recorded. `inbox link` cannot help
there — it needs something to point at. `inbox application <id> <slug>` creates the
application and links the message in one call, dating it by the message rather than
by now(), since the employer replied to something that already existed.

Corrections: `inbox link`/`unlink`, `inbox delete`/`restore`, `inbox read-all`.

## For agents

Pass `--json` for the raw API payload (faithful to the API; ideal for piping):

```bash
rolewatch --json search "site reliability" --limit 5 | jq '.[].public_slug'
rolewatch --json search "golang" --limit 3 | jq -r '.[].description'   # full postings, no extra call
rolewatch --json job <slug> | jq '{title, url}'
```

`search` reads the API's agent endpoint, so every hit already carries the job's
**full description as markdown** — an agent can judge a result set without a
`job <slug>` call per hit. The human table view stays compact (title · company ·
location · slug); the descriptions are in the `--json` payload.

Conventions: results go to **stdout**, errors to **stderr**, and a non-zero exit
code signals failure (e.g. an unauthenticated call exits non-zero with
`run \`rolewatch auth login\``).

## Agent skills and slash commands

Six agent **skills** ship in [`skills/`](skills/), one per task rather than one per
tool, so an agent loads what the question needs and not the whole surface:

| Skill | Covers |
|---|---|
| [`using-rolewatch`](skills/using-rolewatch/SKILL.md) | install, auth, and the map of the rest |
| [`rolewatch-job-search`](skills/rolewatch-job-search/SKILL.md) | profile → facets → search → job/company, contributing a link, reading the inactive-posting signal |
| [`rolewatch-track-applications`](skills/rolewatch-track-applications/SKILL.md) | apply, save, stage, note, `my`, and reporting a posting that never answered |
| [`rolewatch-market-fit`](skills/rolewatch-market-fit/SKILL.md) | scoring a stack against live vacancy demand |
| [`rolewatch-tailor-cv`](skills/rolewatch-tailor-cv/SKILL.md) | reframing a CV toward one vacancy, and the evidence rule |
| [`rolewatch-mail-triage`](skills/rolewatch-mail-triage/SKILL.md) | pushing, judging and linking application mail |

Drop them into a Claude Code (or compatible) skills directory, or install this repo as a
**Claude Code plugin** to get the skills *and* the slash commands wired up automatically:

```
/plugin marketplace add https://github.com/irfanMiral/rolewatch-cli.git
/plugin install rolewatch@rolewatch-cli
```

The plugin adds [`commands/`](commands/): `/job-search`, `/market-fit`,
`/tailor-cv`, `/track-applications` and `/triage-inbox` — each drives its skill's
workflow end to end from one line.

## Configuration

| What | Source (precedence: env → creds file → default) |
|------|--------------------------------------------------|
| Token | `ROLEWATCH_TOKEN` → `~/.rolewatch/creds.json` |
| API base URL | `ROLEWATCH_API_URL` → creds file → `https://rolewatch.irfanmiral.com` |

The token can be supplied entirely via `ROLEWATCH_TOKEN` (no stored file needed),
which suits CI and ephemeral agent sandboxes. `--api-url` overrides the base URL
for a single invocation (e.g. pointing at a local dev server).

## Develop

```bash
go test ./...        # unit tests (config + client + cli), no network
go build ./...
```

## License

MIT — see [LICENSE](LICENSE). The rolewatch backend is MIT too.
