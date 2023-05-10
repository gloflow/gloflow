/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_aws

import (
    "fmt"
    // "context"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/gamelift"
)

//------------------------------------------------------------------------
func GameliftAddPlayer(pPlayerIDstr string) {

	 // Create a new AWS session
	 sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2"), // Replace with your desired region
    })
    if err != nil {
        panic(err)
    }




	svc := gamelift.New(sess)



    input := &gamelift.CreatePlayerSessionInput{
        GameSessionId: aws.String("YOUR_GAME_SESSION_ID"),
        PlayerId:      aws.String("YOUR_PLAYER_ID"),
    }
    
    req, _ := svc.CreatePlayerSessionRequest(input)
    err = req.Send() // context.Background())
    if err != nil {
        fmt.Println("Error creating player session:", err)
        return
    }
    
    // fmt.Println("Player session created successfully:", *resp.PlayerSession.PlayerSessionId)
    

    /*
	// Define player attributes
    playerAttributes := []*gamelift.AttributeValue{
        {
            Key:   aws.String("skill-level"),
            Value: aws.String("beginner"),
        },
    }

    // Create a new player session
    playerSessionInput := &gamelift.CreatePlayerSessionInput{
        GameSessionId: aws.String("GAME_SESSION_ID"), // Replace with the ID of the game session the player is joining
        PlayerId:      aws.String(pPlayerIDstr),

		// custom data you want to associate with the player
        // PlayerData: aws.String("PLAYER_DATA"),
    }

    // Add player attributes to the player session
    playerSessionInput.SetPlayerAttributes(playerAttributes)

    // Create the player session
    playerSessionOutput, err := svc.CreatePlayerSession(playerSessionInput)
    if err != nil {
        panic(err)
    }


    */
}