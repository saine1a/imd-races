module Commands exposing(..)

import Json.Decode exposing (..)
import Model.IMDModel exposing(..)
import Http exposing (..)

athleteDecoder:  Decoder (List Athlete)
athleteDecoder =  
    Json.Decode.list (
        map7 Athlete 
            (field "Athlete" string)
            (field "BirthYear" string)
            (field "OverallRank" int)
            (field "OverallPoints" int)
            (field "SLPointTotal" int)
            (field "GSPointTotal" int)
            (field "SGPointTotal" int)
    )


api : String
api =
    "http://localhost:8080/athlete"


getAthletes : Http.Request (List Athlete)
getAthletes =
    Http.get api ( at ["AllPoints"] athleteDecoder ) 


