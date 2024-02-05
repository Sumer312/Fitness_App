run:
	@templ generate 
	@go build
	@./Health-App-Backend
build:
	@templ generate 
	@go build
