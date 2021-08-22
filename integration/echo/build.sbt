val Http4sVersion = "1.0.0-M21"
val CirceVersion = "0.14.1"
val FinagleVersion = "21.8.0"

libraryDependencies ++= Seq(
  "org.http4s" %% "http4s-blaze-server" % Http4sVersion,
  "org.http4s" %% "http4s-circe" % Http4sVersion,
  "org.http4s" %% "http4s-dsl" % Http4sVersion,
  "io.circe" %% "circe-generic" % CirceVersion,
  "com.twitter" %% "finagle-core" % FinagleVersion
)

assemblyMergeStrategy in assembly := {
  case PathList("javax", "servlet", xs @ _*)        => MergeStrategy.last
  case PathList("javax", "activation", xs @ _*)     => MergeStrategy.last
  case PathList("org", "apache", xs @ _*)           => MergeStrategy.last
  case PathList("com", "google", xs @ _*)           => MergeStrategy.last
  case PathList("com", "esotericsoftware", xs @ _*) => MergeStrategy.last
  case PathList("com", "codahale", xs @ _*)         => MergeStrategy.last
  case PathList("com", "yammer", xs @ _*)           => MergeStrategy.last
  case "module-info.class"                          => MergeStrategy.first
  case PathList("scala", "annotation", "nowarn.class" | "nowarn$.class") =>
    MergeStrategy.first
  case "about.html"                 => MergeStrategy.rename
  case "META-INF/ECLIPSEF.RSA"      => MergeStrategy.last
  case "META-INF/mailcap"           => MergeStrategy.last
  case "META-INF/mimetypes.default" => MergeStrategy.last
  case "plugin.properties"          => MergeStrategy.last
  case "log4j.properties"           => MergeStrategy.last
  case x =>
    val oldStrategy = (assemblyMergeStrategy in assembly).value
    oldStrategy(x)
}
