package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var errVarNotExists error = errors.New("variable does not exist")

type EnvironmentVariables struct {
	ID   string
	vars map[string]string
}

func (e *EnvironmentVariables) ListVariables() []string {
	var variables []string
	for key := range e.vars {
		variables = append(variables, key)
	}
	return variables
}

func (e *EnvironmentVariables) Get(key string) (string, error) {
	val, ok := e.vars[key]
	if !ok {
		return "", errVarNotExists
	}
	return val, nil
}

func (e *EnvironmentVariables) Set(key, value string) {
	e.vars[key] = value
}

type Storage interface {
	Write(string, *EnvironmentVariables) error
	Read(string) (*EnvironmentVariables, error)
}

type Serializer interface {
	Serialize(*EnvironmentVariables) ([]byte, error)
	Deserialize([]byte) (*EnvironmentVariables, error)
}

type PlainTextSerializer struct{}

func (p *PlainTextSerializer) Serialize(env *EnvironmentVariables) ([]byte, error) {
	keys := make([]string, 0, len(env.vars))
	for k, _ := range env.vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var doc string
	for _, key := range keys {
		doc += fmt.Sprintf("%s=%s\n", key, env.vars[key])
	}
	return []byte(doc), nil
}

func (p *PlainTextSerializer) Deserialize(data []byte) (*EnvironmentVariables, error) {
	env := &EnvironmentVariables{
		vars: make(map[string]string),
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line in document: %s", line)
		}
		env.vars[parts[0]] = parts[1]
	}
	return env, nil
}

type DOSpaceStorge struct {
	s3Client   *s3.Client
	spaceName  string
	serializer Serializer
}

func NewDOSpaceStorage(spaceName string, region string, accessKey string, secretKey string) (*DOSpaceStorge, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})),
		config.WithRegion(region),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}
	c := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.digitaloceanspaces.com", region))
	})
	s := &DOSpaceStorge{
		s3Client:   c,
		spaceName:  spaceName,
		serializer: &PlainTextSerializer{},
	}

	return s, nil
}

func (s *DOSpaceStorge) Write(env *EnvironmentVariables) error {
	// format the env var data into a document
	doc, err := s.serializer.Serialize(env)
	if err != nil {
		return fmt.Errorf("unable to serialize environment variables: %v", err)
	}
	// write the document to the space
	_, err = s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.spaceName),
		Key:    aws.String(env.ID),
		Body:   bytes.NewReader(doc),
	})
	if err != nil {
		return fmt.Errorf("unable to write object to space: %v", err)
	}
	return nil
}

func (s *DOSpaceStorge) Read(id string) (*EnvironmentVariables, error) {
	// read the document from the space
	resp, err := s.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.spaceName),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to read object from space: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read object data: %v", err)
	}
	env, err := s.serializer.Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize environment variables: %v", err)
	}
	return env, nil
}
