# orepack
Depend on Sponge Plugins hosted on https://ore.spongepowered.org

## Gradle
```gradle
repositories {
    maven { url 'https://orepack.com' }
}

dependencies {
    compile 'com.orepack.$PLUGIN_AUTHOR:$PLUGIN_ID:$PLUGIN_VERSION'
}
```

## Maven
```xml
<repositories>
  <repository>
    <id>orepack</id>
    <url>https://orepack.com</url>
  </repository>
</repositories>

<dependencies>
  <dependency>
    <groupId>com.orepack.$PLUGIN_AUTHOR</groupId>
    <artifactId>$PLUGIN_ID</artifactId>
    <version>$PLUGIN_VERSION</version>
  </dependency>
</dependencies>
```
