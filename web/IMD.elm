module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)

import Json.Decode exposing (..)

import Http exposing (..)

import Array

type alias Model =
    { athleteName : Maybe String
    }
    
type Msg
    = GotAthletes (Result Http.Error String)
    | GetAthletes
    | NoOp

type alias Athlete = { name : String, class : String }

names = [ Athlete "Rebekah" "U16",Athlete "Caroline" "U16",Athlete "Abi" "U16" ]


decodeContent : Decoder String
decodeContent =
    at [ "AthleteName"] string


athleteToTableRow: Athlete -> Html msg

athleteToTableRow athlete = 
    tr []
    [
        td[][text athlete.name],
        td[][text athlete.class]
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

                Ok name ->
                    let 
                        _ = Debug.log "got name " name
                    in
                        ( { model | athleteName = Just name }, Cmd.none )


api : String
api =
    "http://localhost:8080/athlete/Jonathan%20Brain"


getAthletes : Http.Request String
getAthletes =
    Http.get api decodeContent

view : Model -> Html msg
view model =
    table[class "table table-striped"] (
            List.concat[
            [
                thead[]
                    [ th [][text "Name"]
                    , th [][text "Class"]
                    ]
            ]
            , List.map athleteToTableRow names
         ]
    )

initialization : ( Model, Cmd Msg )
initialization = update GetAthletes { athleteName = Nothing }

main =
    Html.program
    {
        init = initialization
        , update = update
        , view = view
        , subscriptions = always Sub.none
    }
