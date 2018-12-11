package main

import (
    "fmt"
    "time"
    "io"
    "os"
    "os/exec"
)

const (
    MaxFloor = 6
    Length = 11
    Width = 11
    Badge = 1
    Compass = 2
    BuffLowerBoundary = 1
    BuffMidBoundary = 80
    BuffUpperBoundary = 100
    MonsterLowerBoundary = 101
    MonsterUpperBoundary = 400
    DoorLowerBoundary = 401
    DoorUpperBoundary = 500
    StairLowerBoundary = -8
    StairUpperBoundary = -6
    ShopLowerBoundary = 501
    ShopUpperBoundary = 600
    ItemLowerBoundary = 601
    ItemUpperBoundary = 700
)

type Player struct {
    level int
    life_value int
    attack int
    defense int
    gold int
    exp int
    yellow_key int
    blue_key int
    red_key int
    status int
}

type Show interface {
    show()
}

type Buff struct {
    id int
    life_value int
    attack int
    defense int
    yellow_key int
    blue_key int
    red_key int
    gold int
    image byte
}

func (b *Buff) show() {
    fmt.Printf("%c", b.image)
}

func (b *Buff) trigger() {
    fmt.Printf("thanks you help me, i can give you something\n")
    time.Sleep(time.Second)
    fmt.Printf("it maybe fit you\n")
}

type Monster struct {
    id int
    name string
    life_value int
    attack int
    defense int
    gold int
    exp int
    damage int
    image byte
}

func (m *Monster) show() {
    fmt.Printf("%c", m.image)
}

type Door struct {
    id int
    yellow_key int
    blue_key int
    red_key int
    image byte
}

func (d *Door) show() {
    fmt.Printf("%c", d.image)
}

type Stair struct {
    id int
    direct int
    image byte
}

func (s *Stair) show() {
    fmt.Printf("%c", s.image)
}

type Item struct {
    id int
    gold int
    exp int
    life_value int
    attack int
    defense int
    yellow_key int
    blue_key int
    red_key int
    describe string
}

func (item *Item) show() {
}

type Shop struct {
    id int
    list []int
    image byte
}

func (s *Shop) show() {
    fmt.Printf("%c", s.image)
}

type Point struct {
    x, y int
}

var player Player
var cur_floor int
var player_pos Point
var map_info [MaxFloor][Length][Width]int
var cur_map [Length][Width]int
var up_pos [MaxFloor]Point
var down_pos [MaxFloor]Point
var dict map[int]Show

func init_player_info() {
    cur_floor = 0
    player_pos.x = 8
    player_pos.y = 5
    player = Player{}
    player.level = 1
    player.life_value = 10000
    player.attack = 10
    player.defense = 10
    player.yellow_key = 1
    player.blue_key = 1
    player.red_key = 1
    player.status = 1
}

func read_map_file() {
    for i := 0; i<MaxFloor; i++ {
        file_name := fmt.Sprintf("floor_%d.txt", i)
        f, err := os.Open(file_name)
        if err != nil {
            panic(err)
        }

        for j := 0; j < Length; j++ {
            for k := 0; k < Width; k++ {
                fmt.Fscanf(f, "%d", &map_info[i][j][k])
            }
        }
        f.Close()
    }
}

func read_pos_file() {
    f, err := os.Open("up_pos.txt")
    if err != nil {
        panic(err)
    }
    for i := 0; i < MaxFloor; i++ {
        x, y := 0, 0
        fmt.Fscanf(f, "%d %d", &x, &y)
        up_pos[i].x, up_pos[i].y = x, y
    }
    f.Close()

    f, err = os.Open("down_pos.txt")
    if err != nil {
        panic(err)
    }
    for i := 0; i < MaxFloor; i++ {
        x, y := 0, 0
        fmt.Fscanf(f, "%d %d", &x, &y)
        down_pos[i].x, down_pos[i].y = x, y
    }
    f.Close()
}

func read_buff_file() {
    f, err := os.Open("buff.txt")
    defer f.Close()
    if err != nil {
        panic(err)
    }

    for {
        var b Buff
        var s string
        _, err = fmt.Fscanf(f, "%d%d%d%d%d%d%d%d%s", &b.id, &b.life_value, &b.attack, &b.defense, &b.yellow_key, &b.blue_key, &b.red_key, &b.gold, &s)
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }

        b.image = s[0]
        dict[b.id] = &b
    }
}

func read_monster_file() {
    f, err := os.Open("monster.txt")
    defer f.Close()
    if err != nil {
        panic(err)
    }

    for {
        var m Monster
        var s string
        _, err = fmt.Fscanf(f, "%d%s%d%d%d%d%d%d%s", &m.id, &m.name, &m.life_value, &m.attack, &m.defense, &m.gold, &m.exp, &m.damage, &s)
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }

        m.image = s[0]
        dict[m.id] = &m
    }
}

func read_door_file() {
    f, err := os.Open("door.txt")
    defer f.Close()
    if err != nil {
        panic(err)
    }

    for {
        var d Door
        var s string
        _, err = fmt.Fscanf(f, "%d%d%d%d%s", &d.id, &d.yellow_key, &d.blue_key, &d.red_key, &s)
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }

        d.image = s[0]
        dict[d.id] = &d
    }
}

func read_stair_file() {
    f, err := os.Open("stair.txt")
    defer f.Close()
    if err != nil {
        panic(err)
    }

    for {
        var st Stair
        var s string
        _, err = fmt.Fscanf(f, "%d%d%s", &st.id, &st.direct, &s)
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }

        st.image = s[0]
        dict[st.id] = &st
    }
}

func read_item_file() {
    f, err := os.Open("item.txt")
    defer f.Close()
    if err != nil {
        panic(err)
    }

    for {
        var i Item
        _, err = fmt.Fscanf(f, "%d%d%d%d%d%d%d%d%d", &i.id, &i.gold, &i.exp, &i.life_value, &i.attack, &i.defense, &i.yellow_key, &i.blue_key, &i.red_key)

        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }

        var c byte
        var temp []byte
        for {
            _, err = fmt.Fscanf(f, "%c", &c)
            if err == io.EOF || c == '\n' {
                break
            }
            temp = append(temp, c)
        }
        i.describe = string(temp)
        dict[i.id] = &i
    }
}

func read_shop_file() {
    f, err := os.Open("shop.txt")
    defer f.Close()
    if err != nil {
        panic(err)
    }

    for {
        var shop Shop
        var s string
        var arr [3]int
        _, err = fmt.Fscanf(f, "%d%d%d%d%s", &shop.id, &arr[0], &arr[1], &arr[2], &s)
        if err == io.EOF {
            break
        }
        shop.list = append(shop.list, arr[0])
        shop.list = append(shop.list, arr[1])
        shop.list = append(shop.list, arr[2])
        shop.image = s[0]
        dict[shop.id] = &shop
    }
}

func init() {
    dict = make(map[int]Show)
    init_player_info()
    read_map_file()
    read_pos_file()
    read_buff_file()
    read_monster_file()
    read_door_file()
    read_stair_file()
    read_item_file()
    read_shop_file()

    for i := 0; i < Length; i++ {
        for j := 0; j < Width; j++ {
            cur_map[i][j] = map_info[0][i][j]
        }
    }
}

func print() {
    fmt.Printf("cur_floor : %d\n", cur_floor)
    fmt.Printf("Life : %d  Attack : %d  Defense : %d  Gold : %d  Exp : %d\n", player.life_value, player.attack, player.defense, player.gold, player.exp)
    fmt.Printf("Key_Y : %d  Key_B : %d  Key_R : %d\n\n", player.yellow_key, player.blue_key, player.red_key)

    for i := 0; i < Length; i++ {
        for j := 0; j < Width; j++ {
            if cur_map[i][j] == -1 {
                fmt.Printf("*")
                continue
            }
            if i == player_pos.x && j == player_pos.y {
                fmt.Printf("@")
                continue
            }
            if cur_map[i][j] == 0 {
                fmt.Printf(" ")
                continue
            }
            dict[cur_map[i][j]].show()
        }
        fmt.Printf("\n")
    }
}

func change_map(floor int, flag bool) {
    for i := 0; i < Length; i++ {
        for j := 0; j < Width; j++ {
            if map_info[cur_floor][i][j] != cur_map[i][j] {
                map_info[cur_floor][i][j] = 0
            }
        }
    }
    for i := 0; i < Length; i++ {
        for j := 0; j < Width; j++ {
            cur_map[i][j] = map_info[floor][i][j]
        }
    }
    player_pos.x = down_pos[floor].x
    player_pos.y = down_pos[floor].y
    if !flag && cur_floor == floor + 1 {
        player_pos.x = up_pos[floor].x
        player_pos.y = up_pos[floor].y
    } 
    cur_floor = floor
}

func open_door(x, y int) {
    id := cur_map[x][y]
    d := dict[id].(*Door)
    if player.yellow_key < d.yellow_key || player.blue_key < d.blue_key || player.red_key < d.red_key {
        return
    }

    player.yellow_key -= d.yellow_key
    player.blue_key -= d.blue_key
    player.red_key -= d.red_key
    cur_map[x][y] = 0
}

func get_buff(x, y int) {
    id := cur_map[x][y]
    b := dict[id].(*Buff)

    if BuffMidBoundary <= id {
        b.trigger()
    }

    player.life_value += b.life_value
    player.attack += b.attack
    player.defense += b.defense
    player.yellow_key += b.yellow_key
    player.blue_key += b.blue_key
    player.red_key += b.red_key
    player.gold += b.gold
    cur_map[x][y] = 0
}

func cal_damage(m *Monster) int {
    extra := player.life_value * m.damage / 100
    if player.attack <= m.defense {
        return -1
    }
    if player.defense >= m.attack {
        return extra
    }

    temp := m.life_value
    delta1 := player.attack - m.defense
    delta2 := m.attack - player.defense
    if temp % delta1 == 0 {
        return extra + (temp / delta1 - 1) * delta2
    }
    return extra + temp / delta1 * delta2
}

func use_badge() {
    if player.status & Badge == 0 {
        return
    }

    vis := make(map[string]bool)
    for i := 0; i < Length; i++ {
        for j := 0; j < Width; j++ {
            id := cur_map[i][j]
            if id < MonsterLowerBoundary || id > MonsterUpperBoundary {
                continue
            }
            m := dict[id].(*Monster)
            name := m.name
            _, ok := vis[name]
            if ok {
                continue
            }

            vis[name] = true
            damage := cal_damage(m)
            fmt.Printf("name : %s\tattack : %d\tgold/exp : %d-%d\n", m.name, m.attack, m.gold, m.exp)
            fmt.Printf("life : %d\tdefense : %d\t", m.life_value, m.defense)
            if damage != -1 {
                fmt.Printf("damage : %d\n", damage)
            } else {
                fmt.Printf("damage : ???\n")
            }
        }
    }
}

func fight(x, y int) {
    id := cur_map[x][y]
    m := dict[id].(*Monster)
    damage := cal_damage(m)
    if damage == -1 || damage >= player.life_value {
        fmt.Printf("sorry, you cant beat this monster\n")
        return 
    }

    player.life_value -= damage
    player.gold += m.gold
    player.exp += m.exp
    cur_map[x][y] = 0
}

func enter_shop(x, y int) {
    id := cur_map[x][y]
    shop := dict[id].(*Shop)

    fmt.Printf("welcome to shop\n")
    var s string
    for {
        fmt.Printf("Life : %d  Attack : %d  Defense : %d  Gold : %d  Exp : %d\n", player.life_value, player.attack, player.defense, player.gold, player.exp)
        fmt.Printf("Key_Y : %d  Key_B : %d  Key_R : %d\n\n", player.yellow_key, player.blue_key, player.red_key)
        for i, j := range(shop.list) {
            item := dict[j].(*Item)
            fmt.Printf("%d : %s\n", i+1, item.describe)
        }

        fmt.Scanf("%s", &s)
        if s == "q" || s == "Q" {
            return
        }

        if s == "1" || s == "2" || s == "3" {
            index := s[0] - '0' - 1
            item := dict[shop.list[index]].(*Item)

            if player.gold >= item.gold && player.exp >= item.exp {
                player.gold -= item.gold
                player.exp -= item.exp
                player.life_value += item.life_value
                player.attack += item.attack
                player.defense += item.defense
                player.yellow_key += item.yellow_key
                player.blue_key += item.blue_key
                player.red_key += item.red_key
            }
        }        
    }
}

func main(){
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    var s string
    for {
        cmd.Run()
        print()
        fmt.Scanf("%s", &s)
        if s == "q" || s == "Q" {
            break
        }
        if s == "l" || s == "L" {
            use_badge()
            continue
        }

        x, y := player_pos.x, player_pos.y
        if s == "w" || s == "W" {
            x--
        } else if s == "s" || s == "S" {
            x++
        } else if s == "a" || s == "A" {
            y--
        } else if s == "d" || s == "D" {
            y++
        }

        if x < 0 || x >= Length || y < 0 || y >= Width || cur_map[x][y] == -1 {
            continue
        }
        if cur_map[x][y] == 0 {
            player_pos.x, player_pos.y = x, y
            continue
        }

        if StairLowerBoundary <= cur_map[x][y] && cur_map[x][y] <= StairUpperBoundary {
            st := dict[cur_map[x][y]].(*Stair)
            change_map(cur_floor + st.direct, false)
            continue
        }

        if DoorLowerBoundary <= cur_map[x][y] && cur_map[x][y] <= DoorUpperBoundary {
            open_door(x, y)
            continue
        }

        if BuffLowerBoundary <= cur_map[x][y] && cur_map[x][y] <= BuffUpperBoundary {
            get_buff(x, y)
            continue
        }

        if MonsterLowerBoundary <= cur_map[x][y] && cur_map[x][y] <= MonsterUpperBoundary {
            fight(x, y)
            continue
        }

        if ShopLowerBoundary <= cur_map[x][y] && cur_map[x][y] <= ShopUpperBoundary {
            enter_shop(x, y)
            continue
        }
    }
}
