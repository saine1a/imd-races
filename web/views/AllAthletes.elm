module Views.AllAthletes exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)

import Http exposing (..)
import Bootstrap.CDN exposing (..)
import Bootstrap.Table exposing (..)
import Bootstrap.Grid exposing (..)
import Bootstrap.Button exposing (..)
import Html.Events exposing (..)
import Model.IMDModel exposing (..)


import Array


athleteToTableRow: Athlete -> Bootstrap.Table.Row Msg
athleteToTableRow athlete = 
    Bootstrap.Table.tr []
    [
        Bootstrap.Table.td[][text (toString athlete.overallRank)],
        Bootstrap.Table.td[][Bootstrap.Button.button [Bootstrap.Button.primary,Bootstrap.Button.onClick (Model.IMDModel.AthleteClicked athlete.name)][text athlete.name]],
        Bootstrap.Table.td[][text athlete.birthYear],
        Bootstrap.Table.td[][text (toString athlete.overallPoints)],
        Bootstrap.Table.td[][text (toString athlete.slPoints)],
        Bootstrap.Table.td[][text (toString athlete.gsPoints)],
        Bootstrap.Table.td[][text (toString athlete.sgPoints)]
    ]

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
