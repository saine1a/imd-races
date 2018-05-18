module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)

import Json.Decode exposing (..)
import Http exposing (..)
import Bootstrap.CDN exposing (..)
import Bootstrap.Table exposing (..)
import Bootstrap.Grid exposing (..)
import Bootstrap.Button exposing (..)
import Html.Events exposing (..)


import Array

type alias Athlete = { name : String, birthYear : String, overallRank : Int, overallPoints : Int, slPoints : Int, gsPoints : Int, sgPoints : Int}

type alias Model =
    { athletes : List Athlete
    }
    
type Msg
    = GotAthletes (Result Http.Error (List Athlete))
    | GetAthletes
    | NoOp
    | AthleteClicked String

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

athleteToTableRow: Athlete -> Bootstrap.Table.Row Msg
athleteToTableRow athlete = 
    Bootstrap.Table.tr []
    [
        Bootstrap.Table.td[][text (toString athlete.overallRank)],
        Bootstrap.Table.td[][Bootstrap.Button.button [Bootstrap.Button.primary,Bootstrap.Button.onClick (AthleteClicked athlete.name)][text athlete.name]],
        Bootstrap.Table.td[][text athlete.birthYear],
        Bootstrap.Table.td[][text (toString athlete.overallPoints)],
        Bootstrap.Table.td[][text (toString athlete.slPoints)],
        Bootstrap.Table.td[][text (toString athlete.gsPoints)],
        Bootstrap.Table.td[][text (toString athlete.sgPoints)]
    ]


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp -> 
            (model, Cmd.none)

        GetAthletes ->
            let
                _ = Debug.log "Get athletes " 
            in
                ( model, Http.send GotAthletes getAthletes )

        GotAthletes result ->
            case result of
                Err httpError ->
                    let
                        _ =
                            Debug.log "getAthletes error " httpError
                    in
                        ( model, Cmd.none )

                Ok athleteList ->
                    ( { model | athletes = athleteList }, Cmd.none )
        AthleteClicked name ->
            let
                _=
                    Debug.log "Click on " name
            in
                ( model, Cmd.none )


api : String
api =
    "http://localhost:8080/athlete"


getAthletes : Http.Request (List Athlete)
getAthletes =
    Http.get api ( at ["AllPoints"] athleteDecoder ) 

view : Model -> Html Msg

view model =
    Bootstrap.Grid.container[]
        [ 
            Bootstrap.CDN.stylesheet,
                Bootstrap.Table.table
                { options = [ Bootstrap.Table.striped, Bootstrap.Table.hover, Bootstrap.Table.responsive ],
                thead = Bootstrap.Table.simpleThead
                    [
                        Bootstrap.Table.th [][text "Rank"],
                        Bootstrap.Table.th [][text "Name"],
                        Bootstrap.Table.th [][text "Birth Year"],
                        Bootstrap.Table.th [][text "Overall Points"],
                        Bootstrap.Table.th [][text "SL"],
                        Bootstrap.Table.th [][text "GS"],
                        Bootstrap.Table.th [][text "SG"]
                    ]
                    , tbody = Bootstrap.Table.tbody [] (List.map athleteToTableRow model.athletes)
                }
         ]

initialization : ( Model, Cmd Msg )
initialization = update GetAthletes { athletes = [] }

main =
    Html.program
    {
        init = initialization
        , update = update
        , view = view
        , subscriptions = always Sub.none
    }
