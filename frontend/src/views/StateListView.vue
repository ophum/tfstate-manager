<script setup lang="ts">
import { CreateRequest } from '@/gen/api/state/v1/state_pb';
import { client } from '@/main';
import { useAuthStore } from '@/stores/auth';
import { useQuery } from '@tanstack/vue-query';
import { reactive } from 'vue';


const getStates = async () => {
  const res = await client.state.list({})
  return res.states;
}

const auth = useAuthStore()
const { isLoading, isFetching, isError, data, error, refetch } = useQuery({
  queryKey: ['states'],
  queryFn: getStates,
})
if (isError.value) {
  console.log(error)
  auth.sessionExpire()
}

const formData = reactive<CreateRequest>(new CreateRequest);
const create = async () => {
  console.log("on create")
  const res = await client.state.create(formData)
  alert(`created! id: ${res.id} name=${res.name}`)
  formData.name = ""
  formData.description = ""
  refetch();
}
</script>
<template>
  <div class="about">
    <h1>ステート一覧</h1>
    <div v-if="isLoading">Loading...</div>
    <div v-if="isFetching">Fetching...</div>

    <input type="text" v-model="formData.name" placeholder="name" />
    <input type="text" v-model="formData.description" placeholder="description" />
    <button @click="create">create</button>

    <table border="1">
      <thead>
        <tr>
          <th>id</th>
          <th>name</th>
          <th>backend config</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in data" :key="s.id.toString()">
          <td>{{ s.id.toString() }}</td>
          <td>{{ s.name }}</td>
          <td>
            <pre>
  <code>
backend "http" {
  address = "http://localhost:8000/cgi-bin/tfstate-manager/states/{{ s.id }}"
}
  </code>
</pre>

          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
