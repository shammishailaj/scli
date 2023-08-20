/*
Copyright © 2022 Healthians Technology <tech@healthians.com>
*/

package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/shammishailaj/scli/pkg/schemas"
	"github.com/shammishailaj/scli/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

var whitelistAWSCmd = &cobra.Command{
	Use:   "aws",
	Short: "Whitelists Configuration for EC2 and RDS instances",
	Long:  `Whitelists Configuration for EC2 and RDS instances`,
	Run: func(cmd *cobra.Command, args []string) {

		var (
			//accessTypeKeyExists bool
			exitVal      = utils.SUCCEEDED
			awsConfig    aws.Config
			awsConfigErr error
		)

		dryRun, dryRunErr := cmd.Flags().GetBool("dry-run")
		if dryRunErr != nil {
			u.Log.Errorf("Error in doing dry-run %s", dryRunErr.Error())
			u.Log.Infof("Setting dry-run to: false")
			dryRun = false
		}

		awsRegion, awsRegionErr := cmd.Flags().GetString("aws-region")
		if awsRegionErr != nil {
			u.Log.Errorf("Error in getting AWS Access Key ID. %s", awsRegionErr.Error())
			exitVal = utils.AWS_REGION_NOT_PROVIDED
		}

		accessKeyID, accessKeyIDErr := cmd.Flags().GetString("access-key-id")
		if accessKeyIDErr != nil {
			u.Log.Errorf("Error in getting AWS Access Key ID. %s", accessKeyIDErr.Error())
			exitVal = utils.AWS_ACCESS_KEY_NOT_PROVIDED
		}

		secretAccessKey, secretAccessKeyErr := cmd.Flags().GetString("secret-access-key")
		if secretAccessKeyErr != nil {
			u.Log.Errorf("Error in getting AWS Access Key ID. %s", secretAccessKeyErr.Error())
			exitVal = utils.AWS_SECRET_ACCESS_KEY_NOT_PROVIDED
		}

		if accessKeyID == "" || secretAccessKey == "" {
			awsConfig, awsConfigErr = u.AutoAWSConfig(awsRegion)
			if awsConfigErr != nil {
				u.Log.Errorf("Error connecting to AWS. %s", awsConfigErr.Error())
			}
		}

		accessTypeStr, accessTypeStrErr := cmd.Flags().GetString("access-type")
		if accessTypeStrErr != nil {
			u.Log.Errorf("Error in getting Access-type for proxy. %s", accessTypeStrErr.Error())
			u.Log.Infof("Will auto-detect...")
		}

		//accessType, accessTypeKeyExists = accessTypeMap[accessTypeStr] // [How to check if a map contains a key in Go?](https://stackoverflow.com/a/2050629/16898622)
		//if !accessTypeKeyExists {
		//	accessType = accessTypeMap["auto"]
		//}

		accessPort, accessPortErr := cmd.Flags().GetInt64("access-port")
		if accessPortErr != nil {
			u.Log.Errorf("Error in getting accessPort. %s", accessPortErr.Error())
			exitVal = utils.HAPCONFIG_INVALID_BACKEND_PORT
		}

		awsServiceName, awsServiceNameErr := cmd.Flags().GetString("aws-service")
		if awsServiceNameErr != nil {
			u.Log.Errorf("Error in getting AWS Service name. %s", awsServiceNameErr.Error())
			exitVal = utils.AWS_INVALID_SERVICE_NAME
		}

		applicationClassName, applicationClassNameErr := cmd.Flags().GetString("application-class-name")
		if applicationClassNameErr != nil {
			u.Log.Errorf("Error in getting AWS Service name. %s", applicationClassNameErr.Error())
			exitVal = utils.AWS_INVALID_APPLICATION_CLASS_NAME
		}

		portRangeStart, portRangeStartErr := cmd.Flags().GetInt64("port-range-start")
		if portRangeStartErr != nil {
			u.Log.Errorf("Error in getting portRangeStart. %s", portRangeStartErr.Error())
			exitVal = utils.HAPCONFIG_INVALID_PORT_RANGE_START
		}

		portRangeEnd, portRangeEndErr := cmd.Flags().GetInt64("port-range-end")
		if portRangeEndErr != nil {
			u.Log.Errorf("Error in getting portRangeStart. %s", portRangeEndErr.Error())
			exitVal = utils.HAPCONFIG_INVALID_PORT_RANGE_END
		}

		outfile, outfileErr := cmd.Flags().GetString("outfile")
		if outfileErr != nil {
			u.Log.Errorf("Error in getting --outfile. %s", outfileErr.Error())
			exitVal = utils.HAPCONFIG_INVALID_OUTFILE
		}

		if exitVal != utils.SUCCEEDED {
			os.Exit(exitVal)
		}

		switch awsServiceName {
		case "ec2":
			params := &ec2.DescribeInstancesInput{
				DryRun: aws.Bool(dryRun),
				Filters: []types.Filter{
					{
						Name:   aws.String("tag:Name"), // https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
						Values: []string{fmt.Sprintf("%s*", applicationClassName)},
					},
				},
				InstanceIds: nil,
				//MaxResults:  aws.Int32(1000),
				NextToken: nil,
			}
			results, resultsErr := u.GetInstancesDetails(awsConfig, params)
			if resultsErr != nil {
				u.Log.Errorf("Error getting details of instances. %s", resultsErr.Error())
				os.Exit(utils.HAPCONFIG_ERROR_GETTING_INSTANCES_DETAILS)
			}
			u.Log.Infof("Accesstype = %s", accessTypeStr)
			var (
				sshProxySorted                          schemas.HapconfigSSHProxySorted
				sshProxies                              []schemas.HapconfigSSHProxy
				sshProxiesI                             = 0
				generatorHost, generatorHostErr         = os.Hostname()
				generatorUserName, generatorUserNameErr = user.Current()
				outfileData                             = ""
				totalReservations                       = int64(0)
				totalInstances                          = int64(0)
			)

			if generatorHostErr != nil {
				generatorHost = generatorHostErr.Error()
			}

			if generatorUserNameErr != nil {
				generatorUserName.Name = generatorUserNameErr.Error()
			}

			outfileData = fmt.Sprintf("# Generated by %s@%s using obcli generate hapconfig v%s on %s\n# \n\n\n", generatorUserName.Username, generatorHost, Version, time.Now().String())
			for _, reservation := range results.Reservations {
				totalReservations++
				//u.Log.Infof("Reservation ID: %s", *reservation.ReservationId)
				//u.Log.Infof("Instance IDs: ")
				u.Log.Infof("Total instances found for reservationID %s: %d ", *reservation.ReservationId, len(reservation.Instances))
				for _, instance := range reservation.Instances {
					totalInstances++
					//u.Log.Infof(" Instance ID: %s", *instance.InstanceId)
					serverName := ""
					nodeName := ""

					for _, tag := range instance.Tags {
						if *tag.Key == "Name" {
							//u.Log.Infof(" Instance Name: %s", *tag.Value)
							nodeName = *tag.Value
							serverName = strings.ToLower(*tag.Value)
							//break
						}
					}

					if serverName[:len(serverName)-3] != applicationClassName {
						u.Log.Infof("serverName: %s, nodeName := %s, applicationClassName: %s", serverName, serverName[:len(serverName)-3], applicationClassName)
						u.Log.Infoln(" Skipping this server as node name is not the same as application class name provided by the user")
						continue
					}

					// How to get the last X Characters of a Golang String?
					// https://stackoverflow.com/a/26166654/16898622
					nodeNumber, nodeNumberErr := strconv.ParseInt(serverName[len(serverName)-3:], 10, 64)
					if nodeNumberErr != nil {
						u.Log.Fatalf("Error getting node number for server named %s", nodeName)
					}

					if (portRangeStart + nodeNumber - 1) > portRangeEnd {
						u.Log.Fatalf("Unable to assign ports. exiting...")
					}

					bindPort := portRangeStart + nodeNumber - int64(1)

					sshProxies = append(sshProxies, schemas.HapconfigSSHProxy{
						Mode:          "tcp",
						ServerName:    serverName,
						ServerIP:      *instance.PrivateIpAddress,
						ServerPort:    accessPort,
						ServerTimeOut: "2h",
						BindPort:      bindPort,
						ClientTimeOut: "2h",
					})

					//fmt.Printf("\n%s\n\n", sshProxies[sshProxiesI].ProxyString())
					sshProxies[sshProxiesI].ProxyString()
					sshProxiesI++

					//u.Log.Infof(" Instance Type: %s", instance.InstanceType)
					//u.Log.Infof(" Instance Private IP: %s", *instance.PrivateIpAddress)
					//u.Log.Infof(" Instance Platform: %s", instance.Platform.Values())
					//u.Log.Infof(" Instance PlatformDetails: %s", *instance.PlatformDetails)
				}
			}
			u.Log.Infof("Total Reservations found for tag:Name having value matching regular expression %s* - %d", applicationClassName, totalReservations)
			u.Log.Infof("Total Instances found for tag:Name having value matching regular expression %s* - %d", applicationClassName, totalInstances)
			sshProxySorted = sshProxies
			sshProxySorted.Sort()
			for i := 0; i < sshProxiesI; i++ {
				outfileData += sshProxySorted[i].ProxyString() + "\n\n\n"
			}

			u.Log.Infof("Printing proxies in sorted order..........................................................................................\n")
			fmt.Printf("%s\n\n", outfileData)

			bytesWritten, outfileErr := u.OverwriteFile(outfile, outfileData)
			if outfileErr != nil {
				u.Log.Errorf("Error writing to output file %s. %s", outfile, outfileErr.Error())
			} else {
				u.Log.Infof("Successfully wrote %d bytes into file at %s", bytesWritten, outfile)
			}
		}

		if dryRun {
			u.Log.Infof("Dry-run complete")
		}
	},
}

func init() {
	whitelistCmd.AddCommand(whitelistAWSCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whitelistAWSCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	whitelistAWSCmd.Flags().StringP("aws-region", "r", "", "AWS Region")
	whitelistAWSCmd.Flags().StringP("access-key-id", "k", "", "AWS Access Key ID")
	whitelistAWSCmd.Flags().StringP("secret-access-key", "s", "", " AWS Secret Access Key for the Access Key ID provided")
	whitelistAWSCmd.Flags().StringP("access-type", "a", "ssh", "Type of access for which configuration needs to be generated. By default, SSH is generated for Linux/Unix machines and RDP for Windows. Valid values are: ssh, mysql, rdp")
	whitelistAWSCmd.Flags().Int64P("access-port", "b", 22, "Port number using which the service defined by --access-type can be accessed")
	whitelistAWSCmd.Flags().StringP("aws-service", "c", "ec2", "(Optional) AWS service for whose instances the proxy config needs to be generated. Valid values are: ec2, rds. Default: ec2")
	whitelistAWSCmd.Flags().StringP("application-class-name", "d", "", "Name of the application class for which the proxy configuration needs to be generated. IMPORTANT: This command expects that the nodes for a particular application class are named as APPCLASSXXX where APPCLASS is the name and XXX is the node number between 0-999.")
	whitelistAWSCmd.Flags().Int64P("port-range-start", "e", 10000, "Starting port number for the proxy port range for the specified application-class-name")
	whitelistAWSCmd.Flags().Int64P("port-range-end", "f", 20000, "Last port number for the proxy port range for the specified application-class-name")
	whitelistAWSCmd.Flags().StringP("outfile", "g", "", "Path to file where to output the data. File at path gets overwritten")

	//whitelistAWSCmd.Flags().BoolP("dry-run", "r", false, "(Optional) Enable dry-run mode. GRANT queries are only displayed, not executed. Default false")
}
