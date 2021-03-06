package machine

// digitalOceanDescriptions enumerates DigitalOcean instance offerings.
var digitalOceanDescriptions = []Description{
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "ams1", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "ams2", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "ams3", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "blr1", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "fra1", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "lon1", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "nyc1", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "nyc2", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "nyc3", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "sfo1", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "sfo2", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "sgp1", Price: 0.00744},
	{Size: ".5mb", CPU: 1, RAM: .5, Disk: "20", Region: "tor1", Price: 0.00744},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "ams1", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "ams2", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "ams3", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "blr1", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "fra1", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "lon1", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "nyc1", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "nyc2", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "nyc3", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "sfo1", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "sfo2", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "sgp1", Price: 0.01488},
	{Size: "1gb", CPU: 1, RAM: 1, Disk: "30", Region: "tor1", Price: 0.01488},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "ams1", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "ams2", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "ams3", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "blr1", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "fra1", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "lon1", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "nyc1", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "nyc2", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "nyc3", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "sfo1", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "sfo2", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "sgp1", Price: 0.02976},
	{Size: "2gb", CPU: 2, RAM: 2, Disk: "40", Region: "tor1", Price: 0.02976},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "ams1", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "ams2", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "ams3", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "blr1", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "fra1", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "lon1", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "nyc1", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "nyc2", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "nyc3", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "sfo1", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "sfo2", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "sgp1", Price: 0.05952},
	{Size: "4gb", CPU: 2, RAM: 4, Disk: "60", Region: "tor1", Price: 0.05952},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "ams1", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "ams2", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "ams3", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "blr1", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "fra1", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "lon1", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "nyc1", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "nyc2", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "nyc3", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "sfo1", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "sfo2", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "sgp1", Price: 0.11905},
	{Size: "8gb", CPU: 4, RAM: 8, Disk: "80", Region: "tor1", Price: 0.11905},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "ams1", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "ams2", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "ams3", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "blr1", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "fra1", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "lon1", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "nyc1", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "nyc2", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "nyc3", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "sfo1", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "sfo2", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "sgp1", Price: 0.2381},
	{Size: "16gb", CPU: 8, RAM: 16, Disk: "160", Region: "tor1", Price: 0.2381},
	{Size: "m-16gb", CPU: 2, RAM: 16, Disk: "30", Region: "blr1", Price: 0.17857},
	{Size: "m-16gb", CPU: 2, RAM: 16, Disk: "30", Region: "fra1", Price: 0.17857},
	{Size: "m-16gb", CPU: 2, RAM: 16, Disk: "30", Region: "lon1", Price: 0.17857},
	{Size: "m-16gb", CPU: 2, RAM: 16, Disk: "30", Region: "nyc1", Price: 0.17857},
	{Size: "m-16gb", CPU: 2, RAM: 16, Disk: "30", Region: "nyc3", Price: 0.17857},
	{Size: "m-16gb", CPU: 2, RAM: 16, Disk: "30", Region: "sfo2", Price: 0.17857},
	{Size: "m-16gb", CPU: 2, RAM: 16, Disk: "30", Region: "tor1", Price: 0.17857},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "ams2", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "ams3", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "blr1", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "fra1", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "lon1", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "nyc1", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "nyc2", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "nyc3", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "sfo1", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "sfo2", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "sgp1", Price: 0.47619},
	{Size: "32gb", CPU: 12, RAM: 32, Disk: "320", Region: "tor1", Price: 0.47619},
	{Size: "m-32gb", CPU: 4, RAM: 32, Disk: "90", Region: "blr1", Price: 0.35714},
	{Size: "m-32gb", CPU: 4, RAM: 32, Disk: "90", Region: "fra1", Price: 0.35714},
	{Size: "m-32gb", CPU: 4, RAM: 32, Disk: "90", Region: "lon1", Price: 0.35714},
	{Size: "m-32gb", CPU: 4, RAM: 32, Disk: "90", Region: "nyc1", Price: 0.35714},
	{Size: "m-32gb", CPU: 4, RAM: 32, Disk: "90", Region: "nyc3", Price: 0.35714},
	{Size: "m-32gb", CPU: 4, RAM: 32, Disk: "90", Region: "sfo2", Price: 0.35714},
	{Size: "m-32gb", CPU: 4, RAM: 32, Disk: "90", Region: "tor1", Price: 0.35714},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "ams2", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "ams3", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "blr1", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "fra1", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "lon1", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "nyc1", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "nyc2", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "nyc3", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "sfo1", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "sfo2", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "sgp1", Price: 0.71429},
	{Size: "48gb", CPU: 16, RAM: 48, Disk: "480", Region: "tor1", Price: 0.71429},
	{Size: "m-64gb", CPU: 8, RAM: 64, Disk: "200", Region: "blr1", Price: 0.71429},
	{Size: "m-64gb", CPU: 8, RAM: 64, Disk: "200", Region: "fra1", Price: 0.71429},
	{Size: "m-64gb", CPU: 8, RAM: 64, Disk: "200", Region: "lon1", Price: 0.71429},
	{Size: "m-64gb", CPU: 8, RAM: 64, Disk: "200", Region: "nyc1", Price: 0.71429},
	{Size: "m-64gb", CPU: 8, RAM: 64, Disk: "200", Region: "nyc3", Price: 0.71429},
	{Size: "m-64gb", CPU: 8, RAM: 64, Disk: "200", Region: "sfo2", Price: 0.71429},
	{Size: "m-64gb", CPU: 8, RAM: 64, Disk: "200", Region: "tor1", Price: 0.71429},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "ams2", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "ams3", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "blr1", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "fra1", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "lon1", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "nyc1", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "nyc2", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "nyc3", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "sfo1", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "sfo2", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "sgp1", Price: 0.95238},
	{Size: "64gb", CPU: 20, RAM: 64, Disk: "640", Region: "tor1", Price: 0.95238},
	{Size: "m-128gb", CPU: 16, RAM: 128, Disk: "340", Region: "blr1",
		Price: 1.42857},
	{Size: "m-128gb", CPU: 16, RAM: 128, Disk: "340", Region: "fra1",
		Price: 1.42857},
	{Size: "m-128gb", CPU: 16, RAM: 128, Disk: "340", Region: "lon1",
		Price: 1.42857},
	{Size: "m-128gb", CPU: 16, RAM: 128, Disk: "340", Region: "nyc1",
		Price: 1.42857},
	{Size: "m-128gb", CPU: 16, RAM: 128, Disk: "340", Region: "nyc3",
		Price: 1.42857},
	{Size: "m-128gb", CPU: 16, RAM: 128, Disk: "340", Region: "sfo2",
		Price: 1.42857},
	{Size: "m-128gb", CPU: 16, RAM: 128, Disk: "340", Region: "tor1",
		Price: 1.42857},
	{Size: "m-224gb", CPU: 32, RAM: 224, Disk: "500", Region: "blr1", Price: 2.5},
	{Size: "m-224gb", CPU: 32, RAM: 224, Disk: "500", Region: "fra1", Price: 2.5},
	{Size: "m-224gb", CPU: 32, RAM: 224, Disk: "500", Region: "lon1", Price: 2.5},
	{Size: "m-224gb", CPU: 32, RAM: 224, Disk: "500", Region: "nyc1", Price: 2.5},
	{Size: "m-224gb", CPU: 32, RAM: 224, Disk: "500", Region: "nyc3", Price: 2.5},
	{Size: "m-224gb", CPU: 32, RAM: 224, Disk: "500", Region: "sfo2", Price: 2.5},
	{Size: "m-224gb", CPU: 32, RAM: 224, Disk: "500", Region: "tor1", Price: 2.5},
}
