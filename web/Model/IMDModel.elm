module Model.IMDModel exposing (..)

import Http exposing(..)

type alias Athlete = { name : String, birthYear : String, overallRank : Int, overallPoints : Int, slPoints : Int, gsPoints : Int, sgPoints : Int}

type Msg
    = GotAthletes (Result Http.Error (List Athlete))
    | GetAthletes
    | NoOp
    | AthleteClicked String

type alias Model =
    { athletes : List Athlete
    }