package com.echo

import cats._
import cats.effect._
import com.twitter.util.{Return, Throw}
import org.http4s.circe._
import org.http4s._
import io.circe.generic.auto._
import io.circe.syntax._
import org.http4s.dsl._
import org.http4s.implicits._
import org.http4s.server._
import org.http4s.server.blaze.BlazeServerBuilder
import com.twitter.finagle.tracing.TraceId
import org.typelevel.ci.CIStringSyntax
import java.util.Base64

object EchoServer extends IOApp {
  sealed abstract class Reply
  case class Success(data: Trace) extends Reply
  case class Error(error: String) extends Reply

  case class Trace(
      context: String,
      traceId: String,
      spanId: String,
      parent: String,
      sampled: Boolean
  )

  def echoRoutes[F[_]: Monad]: HttpRoutes[F] = {
    val dsl = Http4sDsl[F]
    import dsl._
    HttpRoutes.of[F] { case req @ GET -> Root =>
      val header =
        req.headers.get(ci"L5d-Ctx-Trace").map(_.head.value).getOrElse("")

      TraceId.deserialize(Base64.getDecoder().decode(header)) match {
        case Throw(e) => Ok(Error(e.toString()).asJson)
        case Return(r) => {
          val t = Trace(
            header,
            r.traceIdHigh
              .map(_.toString())
              .getOrElse("0000000000000000") + r.traceId.toString(),
            r.spanId.toString(),
            r.parentId.toString(),
            r.sampled.getOrElse(false)
          )

          Ok(Success(t).asJson)
        }
      }
    }
  }

  import scala.concurrent.ExecutionContext.global

  override def run(args: List[String]): IO[ExitCode] = {
    val apis = Router(
      "/api" -> EchoServer.echoRoutes[IO]
    ).orNotFound

    BlazeServerBuilder[IO](global)
      .bindHttp(5050, "0.0.0.0")
      .withHttpApp(apis)
      .resource
      .use(_ => IO.never)
      .as(ExitCode.Success)
  }
}
