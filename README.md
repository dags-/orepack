# orepack
Plugin dependencies from https://ore.spongepowered.org

## Gradle
```gradle
repositories {
    maven { url 'http://orepack.com' }
}

dependencies {
    compile 'com.orepack:$PLUGIN_ID:$PLUGIN_VERSION'
}
```

## Maven
```xml
<repositories>
  <repository>
    <id>orepack</id>
    <url>http://orepack.com</url>
  </repository>
</repositories>

<dependencies>
  <dependency>
    <groupId>com.orepack</groupId>
    <artifactId>$PLUGIN_ID</artifactId>
    <version>$PLUGIN_VERSION</version>
  </dependency>
</dependencies>
```