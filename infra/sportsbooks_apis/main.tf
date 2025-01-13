module draftkings_basketball_team_api {
    source = "github.com/meagle21/Terraform-Module-Lambda-DynamoDB-Odds-Comparison"
    lambda_function_name = "draftkings_scraping_lambda_nba_team"
    sportsbook_website_url = "https://sportsbook.draftkings.com/leagues/basketball/nba"
    dynambodb_table_name = "draftkings-basketball-team-odds"
    sportsbook = "DraftKings"
    sport = "Basketball"
    team_odds_bool = true
}