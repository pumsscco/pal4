package main

import (
    "fmt"
    "html/template"
    "net/http"
    "github.com/julienschmidt/httprouter"
)
//生成页面
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
    var files []string
    for _, file := range filenames {
        files = append(files, fmt.Sprintf("templates/%s.html", file))
    }
    templates := template.Must(template.ParseFiles(files...))
    w.Header().Set("X-XSS-Protection", "0")
    templates.ExecuteTemplate(w, "layout", data)
}
//首页
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    generateHTML(w, nil, "layout", "navbar", "index")
}
//装备
var routeMapEquip map[string]string=map[string]string{
    "wt":"剑","wl":"双剑","wm":"琴",
    "m":"头部防具","y":"身体防具","x":"足部防具",
    "p":"佩戴",
}
func equipType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    equips:= getEquipType(routeMapEquip[ps.ByName("Type")])
    generateHTML(w, &equips, "layout", "navbar", "equip")
}
//道具
var routeMapProperty map[string][]string=map[string][]string{
    "sw":{"恢复","食物"},"dh":{"恢复","其它恢复类"},
    "dg":{"攻击","攻击类"},
    "cx":{"辅助","香料"},"df":{"辅助","其它辅助类"},
    "ck":{"材料","矿石"},"cs":{"材料","尸块"},"cq":{"材料","其它材料"},
    "jq":{"剧情","剧情类"},
}
func propertyType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    properties:= getPropertyType(routeMapProperty[ps.ByName("Type")])
    generateHTML(w, &properties, "layout", "navbar", "property")
}
//配方
var routeMapPrescription map[string][]string=map[string][]string{
    "wt":{"熔铸图谱","剑"},"wl":{"熔铸图谱","双剑"},"wm":{"熔铸图谱","琴"},
    "m":{"熔铸图谱","头部防具"},"y":{"熔铸图谱","身体防具"},"x":{"熔铸图谱","足部防具"},
    "zz02":{"锻造图谱","锻冶"},
    "zz03":{"注灵图谱","注灵"},
}
func prescriptionType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {   
    prescriptions:= getPrescriptionType(routeMapPrescription[ps.ByName("Type")])
    generateHTML(w, &prescriptions, "layout", "navbar", "prescription")
}
//问答
func questionType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    routeMapQuestion:=map[string]string{"history":"仙剑历史","story":"仙剑故事","world":"仙剑世界"}
    questions:= getQuestionType(routeMapQuestion[ps.ByName("Type")])
    generateHTML(w, &questions, "layout", "navbar", "question")
}
//任务
func missionType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    routeMapMission:=map[string]string{"main":"主线","delegation":"委托","branch":"支线"}
    missions:= getMissionType(routeMapMission[ps.ByName("Type")])
    generateHTML(w, &missions, "layout", "navbar", "mission")
}
//仙术
func magicType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    routeMapMagic:=map[string]string{"self":"我方","enemy":"敌方"}
    magics:= getMagicType(routeMapMagic[ps.ByName("Type")])
    generateHTML(w, &magics, "layout", "navbar", "magic")
}
//特技
func stuntType(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    routeMapStunt:=map[string]string{"self":"我方","enemy":"敌方"}
    stunts:= getStuntType(routeMapStunt[ps.ByName("Type")])
    generateHTML(w, &stunts, "layout", "navbar", "stunt")
}
//升级
func upgradeRole(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    routeMapRole:=map[string]string{"tianhe":"云天河","lingsha":"韩菱纱","mengli":"柳梦璃","ziying":"慕容紫英"}
    upgrades:= getUpgradeRole(routeMapRole[ps.ByName("Name")])
    generateHTML(w, &upgrades, "layout", "navbar", "upgrade")
}
//怪物汇总
func enemySummary(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    enemySummarys:= getEnemySummary()
    //logger.Println("enemys total   summary: ",enemySummarys)
    generateHTML(w, &enemySummarys, "layout", "navbar", "enemy/summary")
}
//怪物明细
func enemyCommon(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    enemyCommons:= getEnemyCommon()
    generateHTML(w, &enemyCommons, "layout", "navbar", "enemy/common")
}
func enemyBasic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    enemyBasics:= getEnemyBasic()
    generateHTML(w, &enemyBasics, "layout", "navbar", "enemy/basic")
}
func enemyResistance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    enemyResistances:= getEnemyResistance()
    generateHTML(w, &enemyResistances, "layout", "navbar", "enemy/resistance")
}
func enemySkill(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    enemySkills:= getEnemySkill()
    generateHTML(w, &enemySkills, "layout", "navbar", "enemy/skill")
}
func enemyDrop(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    enemyDrops:= getEnemyDrop()
    generateHTML(w, &enemyDrops, "layout", "navbar", "enemy/drop")
}
//场景拾取
func pickUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    routeMapPickup:=map[string]string{"m":"迷宫","c":"城镇"}
    sceneType:=routeMapPickup[ps.ByName("Type")]
    if r.Method=="GET" {
        scenes:=getScenes(sceneType)
        generateHTML(w, &scenes, "layout", "navbar", "pickup/scene")
    } else if r.Method=="POST" {
        scene:=r.PostFormValue("scene")
        items:=getPickup(scene)
        generateHTML(w, &items, "layout", "navbar", "pickup/item")
    }
}
//收集详情
func findItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    cat:=ps.ByName("cat")
    t:=ps.ByName("Type")
    if r.Method=="GET" {
        things:=getThings(cat,t)
        generateHTML(w, &things, "layout", "navbar", "find/type")
    } else if r.Method=="POST" {
        thing:=r.PostFormValue("thing")
        methods:=findMethod(cat,t,thing)
        generateHTML(w, &methods, "layout", "navbar", "find/method")
    }
}