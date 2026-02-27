# VHS tape file for chlog demo GIF
Output assets/demo.gif

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

# --- Hidden setup ---
Hide
Type "cd /tmp && rm -rf chlog-demo && mkdir chlog-demo && cd chlog-demo"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 300ms
Show

# --- Demo starts ---

# 1. Initialize
Type "chlog init --project my-app"
Enter
Sleep 1.5s

# 2. Add entries
Type@30ms "chlog add -c added 'User authentication with JWT'"
Enter
Sleep 1s

Type@30ms "chlog add -c fixed 'Memory leak in connection pool'"
Enter
Sleep 1s

# 3. Release
Type "chlog release 1.0.0"
Enter
Sleep 3.5s

Ctrl+L

# 4. Sync to markdown
Type "chlog sync"
Enter
Sleep 1s

# 5. Show in terminal
Type "chlog show"
Enter
Sleep 5s
