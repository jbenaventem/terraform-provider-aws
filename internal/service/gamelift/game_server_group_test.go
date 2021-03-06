package gamelift_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/gamelift"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccGameliftGameServerGroup_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "gamelift", regexp.MustCompile(`gameservergroup/.+`)),
					acctest.MatchResourceAttrRegionalARN(resourceName, "auto_scaling_group_arn", "autoscaling", regexp.MustCompile(`autoScalingGroup:.+`)),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "balancing_strategy", gamelift.BalancingStrategySpotPreferred),
					resource.TestCheckResourceAttr(resourceName, "game_server_protection_policy", gamelift.GameServerProtectionPolicyNoProtection),
					resource.TestCheckResourceAttr(resourceName, "launch_template.0.version", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
		},
	})
}

func TestAccGameliftGameServerGroup_AutoScalingPolicy(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigAutoScalingPolicy(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_policy.0.estimated_instance_warmup", "60"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_policy.0.target_tracking_configuration.0.target_value", "77.7"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
		},
	})
}

func TestAccGameliftGameServerGroup_AutoScalingPolicy_EstimatedInstanceWarmup(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigAutoScalingPolicyEstimatedInstanceWarmup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_policy.0.estimated_instance_warmup", "66"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_policy.0.target_tracking_configuration.0.target_value", "77.7"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
		},
	})
}

func TestAccGameliftGameServerGroup_BalancingStrategy(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigBalancingStrategy(rName, gamelift.BalancingStrategySpotOnly),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "balancing_strategy", gamelift.BalancingStrategySpotOnly),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
		},
	})
}

func TestAccGameliftGameServerGroup_GameServerGroupName(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigGameServerGroupName(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "game_server_group_name", rName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
			{
				Config: testAccGameliftGameServerGroupConfigGameServerGroupName(rName, rName+"-new"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "game_server_group_name", rName+"-new"),
				),
			},
		},
	})
}

func TestAccGameliftGameServerGroup_InstanceDefinition(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigInstanceDefinition(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instance_definition.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
			{
				Config: testAccGameliftGameServerGroupConfigInstanceDefinition(rName, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instance_definition.#", "3"),
				),
			},
		},
	})
}

func TestAccGameliftGameServerGroup_InstanceDefinition_WeightedCapacity(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigInstanceDefinitionWeightedCapacity(rName, "1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instance_definition.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_definition.0.weighted_capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_definition.1.weighted_capacity", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
			{
				Config: testAccGameliftGameServerGroupConfigInstanceDefinitionWeightedCapacity(rName, "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "instance_definition.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_definition.0.weighted_capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "instance_definition.1.weighted_capacity", "2"),
				),
			},
		},
	})
}

func TestAccGameliftGameServerGroup_LaunchTemplate_Id(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigLaunchTemplateId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.id", "aws_launch_template.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "launch_template.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "launch_template.0.version", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
		},
	})
}

func TestAccGameliftGameServerGroup_LaunchTemplate_Name(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigLaunchTemplateName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.id", "aws_launch_template.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "launch_template.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "launch_template.0.version", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
		},
	})
}

func TestAccGameliftGameServerGroup_LaunchTemplate_Version(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigLaunchTemplateVersion(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.id", "aws_launch_template.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "launch_template.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "launch_template.0.version", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
		},
	})
}

func TestAccGameliftGameServerGroup_GameServerProtectionPolicy(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigGameServerProtectionPolicy(rName, gamelift.GameServerProtectionPolicyFullProtection),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "game_server_protection_policy", gamelift.GameServerProtectionPolicyFullProtection),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
		},
	})
}

func TestAccGameliftGameServerGroup_MaxSize(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigMaxSize(rName, "1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "max_size", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
			{
				Config: testAccGameliftGameServerGroupConfigMaxSize(rName, "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "max_size", "2"),
				),
			},
		},
	})
}

func TestAccGameliftGameServerGroup_MinSize(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigMinSize(rName, "1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "min_size", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
			{
				Config: testAccGameliftGameServerGroupConfigMinSize(rName, "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "min_size", "2"),
				),
			},
		},
	})
}

func TestAccGameliftGameServerGroup_RoleArn(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigRoleArn(rName, "test1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					acctest.CheckResourceAttrGlobalARN(resourceName, "role_arn", "iam", fmt.Sprintf(`role/%s-test1`, rName)),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.test1", "arn"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
			{
				Config: testAccGameliftGameServerGroupConfigRoleArn(rName, "test2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
					acctest.CheckResourceAttrGlobalARN(resourceName, "role_arn", "iam", fmt.Sprintf(`role/%s-test2`, rName)),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.test2", "arn"),
				),
			},
		},
	})
}

func TestAccGameliftGameServerGroup_VpcSubnets(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_gamelift_game_server_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(gamelift.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, gamelift.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckGameliftGameServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGameliftGameServerGroupConfigVpcSubnets(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpc_subnets"},
			},
			{
				Config: testAccGameliftGameServerGroupConfigVpcSubnets(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGameliftGameServerGroupExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckGameliftGameServerGroupDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).GameLiftConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_gamelift_game_server_group" {
			continue
		}

		input := gamelift.DescribeGameServerGroupInput{
			GameServerGroupName: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeGameServerGroup(&input)

		if tfawserr.ErrCodeEquals(err, gamelift.ErrCodeNotFoundException) {
			continue
		}

		if err != nil {
			return err
		}

		if output != nil {
			return fmt.Errorf("Gamelift Game Server Group (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGameliftGameServerGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has not set its id", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).GameLiftConn

		input := gamelift.DescribeGameServerGroupInput{
			GameServerGroupName: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeGameServerGroup(&input)

		if err != nil {
			return fmt.Errorf("error reading Gamelift Game Server Group (%s): %w", rs.Primary.ID, err)
		}

		if output == nil {
			return fmt.Errorf("Gamelift Game Server Group (%s) not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccGameliftGameServerGroupIamConfig(rName string, name string) string {
	return fmt.Sprintf(`
data "aws_partition" %[2]q {}

resource "aws_iam_role" %[2]q {
  assume_role_policy = <<-EOF
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "autoscaling.amazonaws.com",
              "gamelift.amazonaws.com"
            ]
          },
          "Action": "sts:AssumeRole"
        }
      ]
    }
  EOF
  name = "%[1]s-%[2]s"
}

resource "aws_iam_role_policy_attachment" %[2]q {
  policy_arn = "arn:${data.aws_partition.%[2]s.partition}:iam::aws:policy/GameLiftGameServerGroupPolicy"
  role       = aws_iam_role.%[2]s.name
}
`, rName, name)
}

func testAccGameliftGameServerGroupLaunchTemplateConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_vpc" "test" {
  default = true
}

data "aws_subnets" "test" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.test.id]
  }
}

resource "aws_launch_template" "test" {
  image_id = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  name     = %[1]q

  network_interfaces {
    subnet_id = data.aws_subnets.test.ids[0]
  }
}
`, rName)
}

func testAccGameliftGameServerGroupInstanceTypeOfferingsConfig() string {
	return `
data "aws_ec2_instance_type_offerings" "available" {
  filter {
    name   = "instance-type"
    values = ["c5a.large", "c5a.2xlarge", "c5.large", "c5.2xlarge", "m4.large", "m4.2xlarge", "m5a.large", "m5a.2xlarge", "m5.large", "m5.2xlarge"]
  }
}
`
}

func testAccGameliftGameServerGroupConfig(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccGameliftGameServerGroupConfigAutoScalingPolicy(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  auto_scaling_policy {
    target_tracking_configuration {
      target_value = 77.7
    }
  }
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccGameliftGameServerGroupConfigAutoScalingPolicyEstimatedInstanceWarmup(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  auto_scaling_policy {
    estimated_instance_warmup = 66
    target_tracking_configuration {
      target_value = 77.7
    }
  }
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccGameliftGameServerGroupConfigBalancingStrategy(rName string, balancingStrategy string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  balancing_strategy     = %[2]q
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, balancingStrategy))
}

func testAccGameliftGameServerGroupConfigGameServerGroupName(rName string, gameServerGroupName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, gameServerGroupName))
}

func testAccGameliftGameServerGroupConfigInstanceDefinition(rName string, count int) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = slice(sort(tolist(data.aws_ec2_instance_type_offerings.available.instance_types)), 0, %[2]d)
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, count))
}

func testAccGameliftGameServerGroupConfigInstanceDefinitionWeightedCapacity(rName string, weightedCapacity string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = slice(sort(tolist(data.aws_ec2_instance_type_offerings.available.instance_types)), 0, 2)
    content {
      instance_type     = instance_definition.value
      weighted_capacity = %[2]q
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, weightedCapacity))
}

func testAccGameliftGameServerGroupConfigLaunchTemplateId(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccGameliftGameServerGroupConfigLaunchTemplateName(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    name = aws_launch_template.test.name
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccGameliftGameServerGroupConfigLaunchTemplateVersion(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id      = aws_launch_template.test.id
    version = 1
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName))
}

func testAccGameliftGameServerGroupConfigMaxSize(rName string, maxSize string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = %[2]s
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, maxSize))
}

func testAccGameliftGameServerGroupConfigMinSize(rName string, minSize string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 2
  min_size = %[2]s
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, minSize))
}

func testAccGameliftGameServerGroupConfigGameServerProtectionPolicy(rName string, gameServerProtectionPolicy string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name        = %[1]q
  game_server_protection_policy = %[2]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.test.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, gameServerProtectionPolicy))
}

func testAccGameliftGameServerGroupConfigRoleArn(rName string, roleArn string) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, roleArn),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size = 1
  min_size = 1
  role_arn = aws_iam_role.%[2]s.arn

  vpc_subnets = [data.aws_subnets.test.ids[0]]

  depends_on = [aws_iam_role_policy_attachment.%[2]s]
}
`, rName, roleArn))
}

func testAccGameliftGameServerGroupConfigVpcSubnets(rName string, count int) string {
	return acctest.ConfigCompose(
		acctest.ConfigLatestAmazonLinuxHvmEbsAmi(),
		testAccGameliftGameServerGroupIamConfig(rName, "test"),
		testAccGameliftGameServerGroupInstanceTypeOfferingsConfig(),
		testAccGameliftGameServerGroupLaunchTemplateConfig(rName),
		fmt.Sprintf(`
resource "aws_gamelift_game_server_group" "test" {
  game_server_group_name = %[1]q
  dynamic "instance_definition" {
    for_each = data.aws_ec2_instance_type_offerings.available.instance_types
    content {
      instance_type = instance_definition.value
    }
  }
  launch_template {
    id = aws_launch_template.test.id
  }

  max_size    = 1
  min_size    = 1
  role_arn    = aws_iam_role.test.arn
  vpc_subnets = slice(data.aws_subnets.test.ids, 0, %[2]d)

  depends_on = [aws_iam_role_policy_attachment.test]
}
`, rName, count))
}
