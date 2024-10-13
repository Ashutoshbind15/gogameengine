package realtime

import (
	"fmt"
	"time"

	"github.com/Ashutoshbind15/gogameengine/internal/types"
	lua "github.com/yuin/gopher-lua"
)

// todo: export the types from the types package

type GameInstanceResources struct {
	Players []types.Player
	Time time.Time
	TimeStarted time.Time
}

// todo: separate out a gameinstance for the dynamic/realtime parts of the game to handle in mem, from the metadata

type Game struct {
	Id string
	InstanceId string
	Gamestate *types.GameState
	Clients []*GameClient
	Moves []types.GameAction // just in case if we want to reconstruct the state from the moves
	GameInfo *GameMeta
	InstanceResources *GameInstanceResources
	Aggregator chan []byte
	TurnBitmap string
}

type RequestPayload struct {
    Script string                 `json:"script"`
	FunctionName string			  `json:"fname"`
    Params map[string]interface{} `json:"params"`
}

func getFunctionSignature(L *lua.LState, fname string) (map[string]string, error) {

	funcsigname := fmt.Sprintln(fname, ":", "signature")
    sig := L.GetGlobal(funcsigname)

    if sig == lua.LNil {
        return nil, fmt.Errorf("signature table not found in script")
    }
    sigTable, ok := sig.(*lua.LTable)
    if !ok {
        return nil, fmt.Errorf("signature is not a table")
    }
    paramsValue := L.GetField(sigTable, "params")
    paramsTable, ok := paramsValue.(*lua.LTable)
    if !ok {
        return nil, fmt.Errorf("params is not a table")
    }

    mp := map[string]string{}
    paramsTable.ForEach(func(k, v lua.LValue) {
        paramTable, ok := v.(*lua.LTable)
        if !ok {
            return
        }
        name := L.GetField(paramTable, "name").String()
        typ := L.GetField(paramTable, "type").String()
		mp[name] = typ
    })
    return mp, nil
}

func goValueToLuaValue(L *lua.LState, val interface{}, expectedType string) (lua.LValue, error) {
    switch expectedType {
    case "int":
        floatVal, ok := val.(float64) // JSON numbers are decoded as float64
        if !ok {
            return nil, fmt.Errorf("expected int, got %T", val)
        }
        intVal := int(floatVal)
        return lua.LNumber(intVal), nil
    case "string":
        strVal, ok := val.(string)
        if !ok {
            return nil, fmt.Errorf("expected string, got %T", val)
        }
        return lua.LString(strVal), nil
    // Add other types as needed
    default:
        return nil, fmt.Errorf("unsupported type %s", expectedType)
    }
}

// Currently each scripts loads, executes and then closes, without interactions with other scripts
func scriptCallerWithParams(payload RequestPayload) (error) {

	L := lua.NewState()
    defer L.Close()

	// register the curr game state

	if err := L.DoString(payload.Script); err != nil {
        return err
    }

	scrparams, err := getFunctionSignature(L, payload.FunctionName)	

	if err != nil {
		return err
	}

	if len(scrparams) != len(payload.Params) {
		// error out
		return fmt.Errorf("CANNOT VALIDATE IP")
	}

	luaArgs := []lua.LValue{}

	for parname, dtype := range scrparams {
        val, ok := payload.Params[parname]
        if !ok {
            return fmt.Errorf("missing parameter %s", parname)
        }
        // Convert val to lua.LValue based on type
        luaVal, err := goValueToLuaValue(L, val, dtype)
        if err != nil {
            return fmt.Errorf("parameter %s: %v", parname, err)
        }
        luaArgs = append(luaArgs, luaVal)
    }

	mutationScriptName := fmt.Sprintln(payload.FunctionName, ":", "modify")

	if err := L.CallByParam(lua.P{
        Fn:      L.GetGlobal(mutationScriptName),
        NRet:    0,
        Protect: true,
    }, luaArgs...); err != nil {
        return err
    }

	// modifiedCstate, err := getGameStateFromLua(L)
    // if err != nil {
    //     return cstate, err
    // }


	return nil
}

func (gm *Game) Runner() {
	for bcastmsg := range gm.Aggregator {
		for _, client := range gm.Clients {
			client.Send <- bcastmsg
		}
	}
}