(generated with Amp)

# Bramble v1 Roadmap

## what's done
- [x] discord bot connects, receives messages, responds
- [x] openrouter integration (single model, stateless HTTP calls)
- [x] agent struct with per-session conversation history
- [x] typing indicator while thinking
- [x] bot ignores itself, responds to DMs + @mentions
- [x] config package stubbed (TOML parsing)

## phase 1: config-driven agent
- [x] wire TOML config into main (load from `~/.bramble/config.toml`)
- [x] move model, system prompt, provider into config
- [x] session descriptions in config (`[agent.sessions]` map)
- [x] inject session description into system prompt per-channel
- [x] create `~/.bramble/` dir + default config on first run

## phase 2: workspace + memory
- [ ] create `~/.bramble/workspace/` on startup
- [ ] `read_file` tool (scoped to workspace)
- [ ] `write_file` tool (scoped to workspace)
- [ ] seed `memory.md` and `pending.md` on first run
- [ ] inject workspace files as context (configurable via `context_files`)

## phase 3: tool calling loop
- [ ] define tool schema format (match openrouter/openai function calling spec)
- [ ] agent loop: model responds → check for tool calls → execute → feed results back → loop until final text response
- [ ] stream progress updates to discord as tools run
- [ ] `read_file` and `write_file` as first tools

## phase 4: more tools
- [ ] `shell` tool (sandboxed, block dangerous commands)
- [ ] `web_search` tool
- [ ] `ssh` tool
- [ ] `switch_model` tool (agent picks its own model mid-conversation)

## phase 5: cron + scheduling
- [ ] `create_cron` tool (agent schedules recurring jobs)
- [ ] `schedule_once` tool (one-time future tasks)
- [ ] per-job model override (cheap model for heartbeats, smart model for reflection)
- [ ] persist schedules to disk, restore on restart
- [ ] default evening reflection cron

## phase 6: slash commands + polish
- [ ] `/reset` — clear session history
- [ ] `/compact` — summarize history into a shorter context
- [ ] `/model` — switch model for current session
- [ ] discord message chunking (2000 char limit)
- [ ] error handling hardening (nil choices, empty responses, rate limits)
- [ ] graceful shutdown (finish in-flight requests)
