import axios from "axios"
const ENDPOINT_URL="http://localhost:8080/api/v1";

const plogAPI={
    async getALL(){
        const result = await axios.get(ENDPOINT_URL);
        return result.data;
    },
    async post(newData){
        const result = await axios.post(ENDPOINT_URL,newData);
        return result.data
    },
    async patch(id,update){
        const result = await axios.patch(`${ENDPOINT_URL}/${id}`,update);
        return result.data;
    }
}

export default plogAPI