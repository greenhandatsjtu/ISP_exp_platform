<template>
    <v-bottom-sheet v-model="sheet" inset max-width="500">
        <template v-slot:activator="{ on, attrs }">
            <v-btn fab fixed bottom right color="primary" v-bind="attrs" v-on="on">
                <v-icon x-large>
                    mdi-plus
                </v-icon>
            </v-btn>
        </template>
        <v-sheet class="text-center pa-4">
            <v-form v-model="valid" ref="form">
                <v-text-field v-model="title"
                              required
                              counter
                              label="标题"
                              :rules="titleRules"></v-text-field>
                <v-text-field v-model="author"
                              counter
                              label="作者（可选，默认为用户名）"
                              :rules="authorRules"></v-text-field>
                <v-textarea v-model="body"
                            required
                            counter
                            label="正文"
                            :rules=" bodyRules"></v-textarea>
                <v-divider/>
                <v-row class="mx-auto">
                    <v-btn
                            color="error"
                            @click="sheet=false"
                    >取消
                    </v-btn>
                    <v-spacer></v-spacer>
                    <v-btn
                            :disabled="!valid"
                            color="info"
                            @click="submit"
                    >确认
                    </v-btn>
                </v-row>
            </v-form>
        </v-sheet>
    </v-bottom-sheet>
</template>

<script>
export default {
    name: "AddNoticeCard",
    inject: ['message','reload'],
    data(){
        return{
            title:'',
            author:'',
            body:'',
            valid: false,
            sheet:false,
            titleRules: [
                v => !!v || "标题不能为空！",
                v => v.length <= 50 || '标题过长！',],
            authorRules: [
                v => v.length <= 50 || '作者名过长！',],
            bodyRules: [
                v => !!v || "正文不能为空！",
                v => v.length <= 200 || '正文过长！',]
        }
    },
    methods:{
        submit(){
            this.$axios.post('admin/notice', {
                author: this.author,
                body:this.body,
                title:this.title})
                .then(() => {
                    this.message("success","添加成功")
                    this.sheet=false
                    this.reload()
                })
                .catch((err) => {
                    console.log(err)
                    this.message("error","添加失败")
                })
        }
    }
}
</script>

<style scoped>

</style>
