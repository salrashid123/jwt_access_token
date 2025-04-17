using System;

using System;
using System.Collections.Generic;
using Google.Cloud.PubSub.V1;
using Google.Protobuf;
using Grpc.Core;
using Grpc.Auth;
using Google.Apis.Auth;
using Google.Apis.Auth.OAuth2;
using Google.Api.Gax.ResourceNames;

namespace main
{
    class Program
    {
        const string projectID = "core-eso";
        [STAThread]
        static void Main(string[] args)
        {
            new Program().Run().Wait();
        }

        private Task Run()
        {
            var _audience = "https://pubsub.googleapis.com/google.pubsub.v1.Publisher";
            var _scopes = "https://www.googleapis.com/auth/cloud-platform";

            var keyFile = "../certs/jwt-access-svc-account.json";
            var stream = new FileStream(keyFile, FileMode.Open, FileAccess.Read);
            ServiceAccountCredential cred = ServiceAccountCredential.FromServiceAccountData(stream).WithUseJwtAccessWithScopes(true);


            var token = cred.GetAccessTokenForRequestAsync(_scopes).Result;

            Console.WriteLine(token);

            var publisher = new PublisherServiceApiClientBuilder { 
                Credential = cred
            }.Build();

            ProjectName projectName = ProjectName.FromProject(projectID);

            foreach (Topic t in publisher.ListTopics(projectName))
                Console.WriteLine(t.Name);
            return Task.CompletedTask;
        }
    }
}

