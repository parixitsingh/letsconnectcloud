STEPS TO RUN

1. Navigate to ../store/cmd and run (path is the location accessible via environment)
    a. go build -o %path%/store.exe main.go (if windows system)
    OR
    b. go build -o %path%<store.<os specific extension if required>> main.go (generic)
2. USE it via following commands
    a. To list all files -->    store ls 
    b. To add files -->         store add filename.txt filename2.txt
    c. To update files -->      store update filename.txt filename2.txt
    d. To remove file -->       store rm filename
    e. To word count -->        store wc
    f. To frequet word -->      store freq-words --limit|-n 10 --order=asc|dsc