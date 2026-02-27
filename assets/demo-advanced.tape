# VHS tape file for chlog advanced demo GIF
Output assets/demo-advanced.gif

Require chlog

# Terminal styling
Set Shell "bash"
Set FontFamily "Noto Sans Mono"
Set FontSize 20
Set Width 1000
Set Height 550
Set Padding 20
Set Framerate 30

# Window decoration
Set WindowBar Colorful
Set WindowBarSize 40
Set Margin 20
Set MarginFill "#282A36"
Set BorderRadius 10

# Theme
Set Theme "Dracula"

# Timing
Set TypingSpeed 40ms
Set PlaybackSpeed 1.5
Set CursorBlink false

# --- Hidden setup: git repo with conventional commits ---
Hide
Type "cd /tmp && rm -rf chlog-adv && mkdir chlog-adv && cd chlog-adv"
Enter
Sleep 500ms
Type "git init -q"
Enter
Sleep 300ms
Type "git commit --allow-empty -m 'feat: OAuth2 login flow' -q"
Enter
Sleep 300ms
Type "git commit --allow-empty -m 'fix: memory leak in pool' -q"
Enter
Sleep 300ms
Type "git commit --allow-empty -m 'feat: rate limiting' -q"
Enter
Sleep 300ms
Type "chlog init --project my-app"
Enter
Sleep 300ms
Type "chlog add -c added 'Internal refactor' --internal"
Enter
Sleep 300ms
Type "chlog release 1.0.0"
Enter
Sleep 300ms
Type "clear"
Enter
Sleep 300ms
Show

# --- Demo starts ---

# 1. Validate
Type "chlog validate"
Enter
Sleep 2.5s

# 2. Scaffold from git commits
Type "chlog scaffold"
Enter
Sleep 4s

Ctrl+L

# 3. Show with internal entries
Type "chlog show 1.0.0 --internal"
Enter
Sleep 4.5s

Ctrl+L

# 4. Sync --split (public + internal changelogs)
Type "chlog sync --split"
Enter
Sleep 3s

# 6. Check --split
Type "chlog check --split"
Enter
Sleep 3s

Ctrl+L

# 7. Extract as markdown
Type "chlog extract 1.0.0"
Enter
Sleep 5s
