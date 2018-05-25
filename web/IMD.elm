module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Http exposing (..)

import Commands exposing(..)
import Views.AllAthletes exposing(..)

import Model.IMDModel exposing(..)


import Array

update : Msg -> Model.IMDModel.Model -> ( Model.IMDModel.Model, Cmd Msg )

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

initialization : ( Model.IMDModel.Model, Cmd Msg )
initialization = update Model.IMDModel.GetAthletes { athletes = [] }

main =
    Html.program
    {
        init = initialization
        , update = update
        , view = Views.AllAthletes.view
        , subscriptions = always Sub.none
    }
