# Asky Demo Tape

# --- Configure the terminal -------------------------------
Output assets/asky-showcase.gif
Set FontSize 18
Set WindowBar Colorful
Set Margin 20
Set MarginFill "#9E757E"
Set BorderRadius 10
Set Theme "Catppuccin Mocha"

# --- Change prompt symbol ---------------------------------
Hide
Type "PS1='\033[33m$\033[0m '"
Enter
Ctrl+L
Show

# --- Prep: Run the showcase binary ------------------------
Type "./bin/asky-showcase"
Sleep 1s
Enter
Sleep 2s

# --- Prompt Showcase: Text Input --------------------------
Type@150ms "Emplois"
Sleep 1s
Enter
Sleep 2s

# --- Prompt Showcase: Text Input --------------------------
Type@150ms "Musswashans"
Sleep 1s
Enter
Sleep 2s

# --- Prompt Showcase: Number Input ------------------------
Type@150ms "69"
Sleep 1s
Enter
Sleep 2s

# --- Prompt Showcase: Secure Input ------------------------
Type@150ms "max"
Sleep 1s
Enter
Sleep 2s
Backspace 3
Sleep 1s
Type@150ms "Imumoccupancy120"
Sleep 2s
Enter
Sleep 2s

# --- Prompt Showcase: Select ------------------------------
Down
Sleep 0.5s
Down
Sleep 0.5s
Down
Sleep 0.5s
Space
Sleep 0.5s
Down
Sleep 0.5s
Down
Sleep 0.5s
Down
Sleep 0.5s
Down
Sleep 0.5s
Down
Sleep 0.5s
Up
Sleep 0.5s
Up
Sleep 0.5s
Space
Sleep 2s
Tab
Sleep 2.5s
Type@150ms "Java"
Sleep 3s
Down
Sleep 0.5s
Backspace 4
Type@150ms "Erlang"
Sleep 2s
Space
Sleep 2s
Enter
Sleep 2s

# --- Prompt Showcase: Confirmation ------------------------
Type "y"
Sleep 2s
Enter

# --- Indicator: Loading Spinner ---------------------------
Sleep 4s

# --- Indicator: Progress Bar ------------------------------
Sleep 6s

# --- End ---
Sleep 2s