import fs from 'fs';

const context = (notebook, paragraph) => {
  // see ../linux/setup_alpine.sh
  const parent = `/home/notebook/${notebook}_${paragraph}`;
  return {
    // to get saved data from machine
    /**
     * @param {string} key
     * @returns {string}
    */
    get: (key) => {
      let data = fs.readFileSync(`${parent}/${key}`, { encoding: 'utf8', flag: 'r' });
      return data || "";
    },
    // to persist data, value must be string
    /**
     * 
     * @param {string} key 
     * @param {string} value 
     */
    set: (key, value) => {
      if(!fs.existsSync(parent)) {
        fs.mkdirSync(parent);
      }
      fs.writeFileSync(`${parent}/${key}`, value);{}
    },
    /**
     * to write data back to notebook
     * @param {string} data 
     * @param {'json'|'md'|'string'} type 
     */
    put: async (data, type) => {
      // this is host machine Ideally should be somekind of MQ in there. but I am lazy
      const agent_url = `http://localhost:1323/notebook/${notebook}/paragraph/${paragraph}/result`
      // I dont do validate now
      await fetch(agent_url, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ data, type }),
      })
    }
  }
}

function main () {
  // to be filled through go
  let Nb = Object.freeze(context(%s, %s));
  // code goes here
  %s
}

main()