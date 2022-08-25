loc="$(pwd)"

printf "Starting Server in New Window"
osascript -  "$loc"  <<EOF
    on run argv
        tell application "Terminal"
            do script ("cd " & quoted form of item 1 of argv) & " && go run main.go"
        end tell
    end run
EOF

set -x

sleep 3
printf "Sending Website List \n"
curl -X POST localhost:8080/websites -d '{"websites": ["google.com", "yahoo.com", "abcd.com"]}'
sleep 2
printf "Getting status all Websites\n"
curl "localhost:8080/websites"
sleep 2
printf "Getting status of selected Websites\n"
curl 'localhost:8080/websites?name=google.com&name=abcd.com'
sleep 2
printf "Requesting unsupported Website\n"
curl 'localhost:8080/websites?name=facebook.com'
sleep 2
printf "Updating website list\n"
curl -X POST localhost:8080/websites -d '{"websites": ["google.com", "abcd.com"]}'
sleep 2
printf "Trying Previously Supported Website\n"
curl 'localhost:8080/websites?name=yahoo.com'