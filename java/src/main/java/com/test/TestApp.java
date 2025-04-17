package com.test;

import com.google.cloud.pubsub.v1.TopicAdminClient;
import com.google.cloud.pubsub.v1.TopicAdminClient.ListTopicsPagedResponse;
import com.google.cloud.pubsub.v1.TopicAdminSettings;

import com.google.pubsub.v1.ListTopicsRequest;
import com.google.pubsub.v1.ProjectName;
import com.google.pubsub.v1.Topic;
import com.google.api.gax.core.FixedCredentialsProvider;
import com.google.auth.oauth2.GoogleCredentials;
import com.google.auth.oauth2.ServiceAccountCredentials;
import com.google.auth.oauth2.ServiceAccountJwtAccessCredentials;
import java.io.File;
import java.io.InputStream;
import java.net.URI;
import java.io.FileInputStream;

public class TestApp {
	public static void main(String[] args) {
		TestApp tc = new TestApp();
	}

	public TestApp() {
		try {

			String serviceAccuntKey = "../certs/jwt-access-svc-account.json";
			InputStream credentialsStream = new FileInputStream(new File(serviceAccuntKey));

			// URI audience = new URI("https://pubsub.googleapis.com/google.pubsub.v1.Publisher");
			// ServiceAccountJwtAccessCredentials credentials = ServiceAccountJwtAccessCredentials
			// 		.fromStream(credentialsStream, audience);
					
			ServiceAccountCredentials credentials = ServiceAccountCredentials.fromStream(credentialsStream).createWithUseJwtAccessWithScope(true);

			TopicAdminClient topicClient = TopicAdminClient
					.create(TopicAdminSettings.newBuilder()
							.setCredentialsProvider(FixedCredentialsProvider.create(credentials)).build());

			ListTopicsRequest listTopicsRequest = ListTopicsRequest.newBuilder()
					.setProject(ProjectName.format("core-eso"))
					.build();

			ListTopicsPagedResponse response = topicClient.listTopics(listTopicsRequest);
			Iterable<Topic> topics = response.iterateAll();
			for (Topic topic : topics)
				System.out.println(topic);

		} catch (Exception ex) {
			System.out.println("Error: " + ex);
		}
	}

}
