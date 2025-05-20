package firebase

// func InitFirebaseMessaging() *messaging.Client {
// 	workingDir, err := os.Getwd()
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Error getting current working directory")
// 	}

// 	credentialFilePath := filepath.Join(workingDir, config.Envs.FirebaseMessaging.CredentialServiceAccount)

// 	// cfg := &firebase.Config{
// 	// 	ProjectID: config.Envs.FirebaseMessaging.ProjectID,
// 	// }

// 	opt := option.WithCredentialsFile(credentialFilePath)
// 	app, err := firebase.NewApp(context.Background(), &firebase.Config{ProjectID: "fcm-notification-e91d6"}, opt)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("Error connecting to Firebase")
// 	}

// 	fmt.Println("APP : ", app)

// 	client, err := app.Messaging(context.Background())
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to get firebase messaging client")
// 	}

// 	fmt.Println("Client : ", client)

// 	return client
// }
