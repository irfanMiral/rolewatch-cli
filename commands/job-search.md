---
description: Search live vacancies on rolewatch and shortlist the ones worth applying to
argument-hint: [what they are looking for, in their own words]
allowed-tools: Bash(rolewatch:*)
---

Find jobs on rolewatch for: **$ARGUMENTS**

Follow the `rolewatch-job-search` skill. In short:

1. Run `rolewatch --json profile` first. The person has already saved the roles,
   skills, geography and work format they want — search on that, say what you
   searched on, and ask only about what the profile does not answer. If the
   request above is empty, the profile *is* the request.
2. Run `rolewatch --json facets` (scoped with `--category` when the role is
   known) before choosing any filter value, so every value and skill slug is one
   the market actually uses.
3. Search with keyword plus facet filters, `--limit 20`. The JSON carries each
   hit's full description, so screen the set in that one call rather than
   fetching postings one by one.
4. Shortlist. For each job you put forward give the title, company, location,
   salary if stated, and one line on why it fits *this* person. Link it as
   `https://rolewatch.irfanmiral.com/jobs/<public_slug>`, one per line.
5. If a posting carries a `ghost` object, relay it with its hedge and its scale
   ("2 of 4 checks fired") and never as "this job is fake".

If they name a company or paste a link rolewatch has no results for, run
`rolewatch contribute <url>` rather than reporting a dead end.

End by asking whether to mark any of them applied or saved.
