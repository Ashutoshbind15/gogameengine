package extrascripts

var FrontalAttackScript = `
function attack(gs)
    local playersInGame = gs.players

    -- If no players are in the game, return an empty state
    if #playersInGame < 1 then
        return {}
    end

    -- Retrieve the first player's position and the next level
    local player = playersInGame[1]
    local py = player.py
    local nextlvl = py + 1

    -- Get the enemies in the game
    local towersInGame = gs.enemies
    local resEnemies = {}

    -- Loop through all enemies and adjust health if their position matches nextlvl
    for i = 1, #towersInGame do
        local tower = towersInGame[i]
        local ty = tower.py

        -- Only reduce health if the enemy is at py + 1
        if ty == nextlvl then
            table.insert(resEnemies, {
                health = tower.health - 1,  -- Decrease health by 1
                px = tower.px,
                py = tower.py
            })
        else
            -- If enemy is not in the next level, preserve its original state
            table.insert(resEnemies, {
                health = tower.health,  -- Preserve original health
                px = tower.px,
                py = tower.py
            })
        end
    end

    -- Return the updated game state
    return {
        enemies = resEnemies,
        time = gs.time,
        timeleft = gs.timeleft,
        arena = gs.arena,
        players = gs.players
    }
end
`