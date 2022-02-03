
# Init env for tailwindcss

## Pre requsites

Get node

```
curl -sL https://deb.nodesource.com/setup_14.x | sudo bash -
```

Get npm

```
sudo apt install npm
```

## Init project 

```
npm init
```

Update package.json : 


```
{
  "name": "vue1",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "build": "tailwindcss build -i input.css -o output.css"
  },
  "author": "",
  "license": "ISC"
}

```

### Create input.css

```
@tailwind base;
@tailwind components;
@tailwind utilities;
```

Install tailwind :

```
npm install -D tailwindcss@latest postcss@latest autoprefixer@latest
npm run build
```

Add node_modules to the .gitignore file.