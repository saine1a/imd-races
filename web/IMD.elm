module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)

import Json.Decode exposing (..)
import Http exposing (..)

import Array

type alias Athlete = { name : String, birthYear : String }

type alias Model =
    { athletes : List Athlete
    }
    
type Msg
    = GotAthletes (Result Http.Error (List Athlete))
    | GetAthletes
    | NoOp

athleteDecoder:  Decoder (List Athlete)
athleteDecoder =  
    Json.Decode.list (
        map2 Athlete 
            (field "Athlete" string)
            (field "BirthYear" string)
    )

athleteToTableRow: Athlete -> Html msg
athleteToTableRow athlete = 
    tr []
    [
        td[][text athlete.name],
        td[][text athlete.birthYear]
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


api : String
api =
    "http://localhost:8080/athlete"


getAthletes : Http.Request (List Athlete)
getAthletes =
    Http.get api ( at ["AllPoints"] athleteDecoder ) 

view : Model -> Html msg

view model =
    table[class "table table-striped"] (
            List.concat[
            [
                thead[]
                    [ th [][text "Name"]
                    , th [][text "Birth Year"]
                    ]
            ]
            , List.map athleteToTableRow model.athletes
         ]
    )

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
