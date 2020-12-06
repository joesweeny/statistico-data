# gRPC
Statistico's internal systems communicate via gRPC. This application's gRPC specifications can be found in the 
[statistico-proto](https://github.com/statistico/statistico-proto/data) repository. For more on gRPC view
 [here](https://grpc.io/docs/guides/)

This application exposes four services:
- FixtureService
- PlayerStatsService
- ResultService
- TeamStatsService

The parameters required to access these services are well defined in their respective `.proto` files. 

To access this applications services using a local client we recommend [gRPCurl](https://github.com/fullstorydev/grpcurl). 
Example calls are:

#### To fetch fixtures between a date period
```proto
grpcurl \
    -plaintext \
    -d \
    '{"season_id": 16036, "date_from": "2019-04-03T00:00:00+00:00", "date_to": "2019-04-03T23:59:59+00:00"}' \
    localhost:50051  \
    statisticoproto.FixtureService/ListSeasonFixtures
```
#### To fetch a fixture by ID
```proto
grpcurl \
    -plaintext \
    -d \
    '{"fixture_id": 5601}' \
    localhost:50051  \
    statisticoproto.FixtureService/FixtureByID
```

#### To fetch player stats for a given fixture
```proto
grpcurl \
    -plaintext \
    -d \
    '{"fixture_id": 7019}' \
    localhost:50051  \
    statisticoproto.PlayerStatsService/GetPlayerStatsForFixture
```  
    
#### To fetch results for a given Team
```proto
grpcurl \
    -plaintext \
    -d \
    '{"team_id": 501, "limit": 75, "date_before": "2019-04-03T23:59:59+00:00"}' \
    localhost:50051  \
    statisticoproto.ResultService/GetResultsForTeam
```
    
#### To fetch team stats for a given fixture
```proto
grpcurl \
    -plaintext \
    -d \
    '{"fixture_id": 7019}' \
    localhost:50051  \
    statisticoproto.TeamStatsService/GetTeamStatsForFixture
```