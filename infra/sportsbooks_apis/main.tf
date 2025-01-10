variable module_source_github_url {
    description = "Path of the source GitHub repo that helps with deployment of Lambda, DynamoDB and necessary roles, etc."
    default = "github.com/meagle21/Terraform-Module-Lambda-DynamoDB-Odds-Comparison"
}

module draftkings_api {
    source = var.module_source_github_url
    lambda_function_name = "draftkings_scraping_lambda_nba_team"
    sportsbook_website_url = "https://sportsbook.draftkings.com/leagues/basketball/nba"
}