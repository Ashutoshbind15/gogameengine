package extrascripts

var FrontalAttackScript = `
function playerSingleLevelAttackInit(gs)
	local playersInGame = gs.players

	if #playersInGame < 1 then
		return {}
	end

	local player = playersInGame[1]
	local py = player.py
	local nextlvl = py + 1

	local towersInGame = gs.enemies
	local resEnemies = {}

	for i = 1, #towersInGame do
		local tower = towersInGame[i]
		local ty = tower.py
		if ty == nextlvl then
			local chealth = tower.health
			table.insert(resEnemies, {
				health = chealth - 1,
				px = tower.px,
				py = tower.py
			})
		end
	end

	return {
		enemies = resEnemies,
		time = gs.time,
		timeleft = gs.timeleft,
		arena = gs.arena,
		players = gs.players
	}
end
`