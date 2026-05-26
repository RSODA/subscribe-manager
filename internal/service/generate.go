package service

//go:generate sh -c "rm -rf mocks && mkdir mocks"
//go:generate minimock -i SubscriptionService -o ./mocks/ -s "_minimock.go"
