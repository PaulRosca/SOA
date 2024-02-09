import amqp from "amqplib";
import nodemailer from "nodemailer";
import nodemailerSendgrid from "nodemailer-sendgrid";
import 'dotenv/config'

const transporter = nodemailer.createTransport(
nodemailerSendgrid({
     apiKey: process.env.apiKey || ""
  })
);

const queue = "emails";

(async () => {
  try {
    const connection = await amqp.connect("amqp://guest:guest@rabbitmq.default.svc.cluster.local:5672/");
    const channel = await connection.createChannel();

    process.once("SIGINT", async () => {
      await channel.close();
      await connection.close();
    });

    await channel.assertQueue(queue, { durable: false });
    await channel.consume(
      queue,
      (message: any) => {
        if (message) {
          const parsedMessage = JSON.parse(message.content.toString())
          console.log(
            " [x] Received '%s'",
            parsedMessage
          );
          const mailOptions = {
            from: 'paul.rosca@stud.ubbcluj.ro',
            to: parsedMessage.email,
            subject: `Order ${parsedMessage.id}`,
            text: `Order placed successfully, you can see it in the app!`
          };

          transporter.sendMail(mailOptions, function(error, info){
            if (error) {
              console.dir(error, { depth: null });
            } else {
              console.log('Email sent: ' + info.response);
            }
          });
        }
      },
      { noAck: true }
    );

    console.log(" [*] Waiting for messages. To exit press CTRL+C");
  } catch (err) {
    console.warn(err);
  }
})();
