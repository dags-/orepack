window.addEventListener("load", function() {
    let hash = window.location.hash;
    if (hash === "") {
        return;
    }
    document.getElementById("search").value = hash.substring(1);
    find();
});

function get(url) {
    return fetch(`https://cors-anywhere.herokuapp.com/${url}`).catch(console.error);
}

function clear(el) {
    while (el.lastChild) {
        el.removeChild(el.lastChild);
    }
}

function render(html, type) {
    if (!type) {
        type = "div";
    }
    let el = document.createElement(type);
    el.innerHTML = html;
    return el;
}

function splitPath(path) {
    let find = "spongepowered.org";
    let i = path.indexOf(find);
    if (i > 0) {
        i += find.length;
    }
    path = path.substring(i + 1);
    return path.split("/");
}

function setLoading(loading) {
    let content = document.getElementById("content");
    if (loading) {
        let project = document.getElementById("project");
        content.insertBefore(render(renderLoader()).firstChild, content.firstChild);
        clear(project);
    } else {
        let loader = document.getElementById("loader");
        content.removeChild(loader);
    }
}

function goTo(id) {
    document.getElementById(id).scrollIntoView({behavior: "smooth"});
}

function find() {
    let input = document.getElementById("search").value;
    let parts = splitPath(input);
    if (parts.length < 2) {
        return;
    }

    let owner = parts[0];
    let name = parts[1];
    let root = document.getElementById("project");
    setLoading(true);

    getPluginId(owner, name)
        .then(getVersions)
        .then(versions => {
            window.location.href = `#${owner}/${name}`;
            let project = renderProject(owner, name, versions);
            for (let i = 0; i < project.length; i++) {
                root.appendChild(project[i]);
            }
            goTo("search");
        })
        .catch(() => {
            root.innerHTML = renderDescription()
        })
        .finally(() => setLoading(false));
}

function setVersion(owner, id, version) {
    let gradle = document.getElementById("gradle");
    if (gradle) {
        gradle.innerText = renderGradle(owner, id, version);
    }

    let maven = document.getElementById("maven");
    if (maven) {
        maven.innerText = renderMaven(owner, id, version);
    }
}

function getPluginId(owner, project) {
    return get(`https://ore.spongepowered.org/api/v1/users/${owner}`)
        .then(r => r.json())
        .then(user => user["projects"].find(p => p["name"] === project))
        .then(proj => proj["pluginId"]);
}

function getVersions(project) {
    return get(`https://ore.spongepowered.org/api/v1/projects/${project}/versions`)
        .then(r => r.json());
}

function getSpongeDep(version) {
    let sponge = version["dependencies"].find(dep => dep["pluginId"] === "spongeapi");
    if (sponge) {
        return sponge["version"];
    }
    return "unknown";
}

function renderProject(owner, project, versions) {
    let children = [];
    children.push(renderTitle(owner, project));
    children.push(renderVersions(owner, project, versions));
    children.push(renderGradleDependency(owner, versions[0]["pluginId"], versions[0]["name"]));
    children.push(renderMavenDependency(owner, versions[0]["pluginId"], versions[0]["name"]));
    return children;
}

function renderTitle(owner, project) {
   return render(`<div class="project-title">${owner}/${project}</div>`).firstChild;
}

function renderVersions(owner, project, versions) {
    let root = render(`<div class="project-versions">Versions:</div>`).firstChild;
    let table = render(`<table><tr><th>Plugin Version</th><th>Sponge API</th><th>Build</th></tr></table>`).firstChild;
    let tbody = table.firstChild;
    versions.forEach(version => tbody.appendChild(renderVersion(owner, version)));
    root.appendChild(table);
    return root;
}

function renderVersion(owner, version) {
    let ver = `<td>${version["name"]}</td>`;
    let spn = `<td>${getSpongeDep(version)}</td>`;
    let get = `<td onclick="setVersion('${owner}','${version["pluginId"]}','${version["name"]}')"><a onclick="goTo('get')">Get</a></td>`;
    return render(ver + spn + get, "table").firstChild.firstChild;
}

function renderGradleDependency(owner, pluginId, version) {
    let root = render(`<div class="project-dependency">Gradle:</div>`).firstChild;
    let pre = render(`<pre></pre>`).firstChild;
    let code = render(`<code id="get"></code>`).firstChild;
    code.innerText = renderGradle(owner, pluginId, version);
    pre.appendChild(code);
    root.appendChild(pre);
    return root;
}

function renderMavenDependency(owner, pluginId, version) {
    let root = render(`<div class="project-dependency">Maven:</div>`).firstChild;
    let pre = render(`<pre></pre>`).firstChild;
    let code = render(`<code></code>`).firstChild;
    code.innerText = renderMaven(owner, pluginId, version);
    pre.appendChild(code);
    root.appendChild(pre);
    return root;
}

function renderGradle(owner, pluginId, version) {
    return `
repositories {
    maven { url "https://orepack.com" }
}

dependencies {
    compile "com.orepack.${owner}:${pluginId}:${version}"
}
`.trim();
}

function renderMaven(owner, pluginId, version) {
    return `
<repositories>
  <repository>
    <id>orepack</id>
    <url>https://orepack.com</url>
  </repository>
</repositories>

<dependencies>
  <dependency>
    <groupId>com.orepack.${owner}</groupId>
    <artifactId>${pluginId}</artifactId>
    <version>${version}</version>
  </dependency>
</dependencies>
`.trim();
}

function renderLoader() {
    return `<div class="loader" id="loader"></div>`;
}

function renderDescription() {
    return `<div class="project-title">Hello world</div>
        <div>Depend on Sponge plugins hosted on <a href="https://ore.spongepowered.org" target="_blank">Ore</a></div>
        <div>Repository: <code>https://orepack.com</code></div>
        <div>GroupId: <code>com.orepack.$PluginAuthor</code></div>
        <div>ArtifactId: <code>$PluginId</code></div>
        <div>Version: <code>$PluginVersion</code></div>`;
}