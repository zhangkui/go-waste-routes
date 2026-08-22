<template>
  <el-card>
    <p>{{ auth.user?.display_name }}</p>
    <el-form @submit.prevent="changePassword">
      <el-form-item><el-input v-model="form.old_password" type="password" placeholder="旧密码" /></el-form-item>
      <el-form-item><el-input v-model="form.new_password" type="password" placeholder="新密码" /></el-form-item>
      <el-button type="primary" native-type="submit">修改密码</el-button>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import http from '@/api/http'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const form = reactive({ old_password: '', new_password: '' })

async function changePassword() {
  await http.put('/users/password', form)
}
</script>
