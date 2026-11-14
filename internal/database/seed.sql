INSERT OR IGNORE INTO employer(id, name) VALUES
  (1, 'Feed the Children'),
  (2, 'MAMALAND (A MOTA-ENGIL GROUP COMPANY)'),
  (3, 'Africa Global Logistics Malawi Limited'),
  (4, 'Trident Security Company Limited'),
  (5, 'Strategic Sales and Merchandising Services Limited');

INSERT OR IGNORE INTO skill(name, category) VALUES
  ('DHIS2', 'data platform'),
  ('KoboToolbox', 'data platform'),
  ('CommCare', 'data platform'),
  ('Power BI', 'data platform'),
  ('Tableau', 'data platform'),
  ('Google Earth Engine', 'data platform'),
  ('Survey123', 'data platform'),
  ('GIS', 'gis'),
  ('QGIS', 'gis'),
  ('ArcGIS', 'gis'),
  ('ArcGIS Pro', 'gis'),
  ('ArcGIS Online', 'gis'),
  ('Remote sensing', 'gis'),
  ('Spatial analysis', 'gis'),
  ('Cartography', 'gis'),
  ('Satellite imagery', 'gis'),
  ('LiDAR', 'gis'),
  ('Drone operations', 'gis'),
  ('MRV', 'gis'),
  ('Carbon projects', 'gis'),
  ('Biomass estimation', 'gis'),
  ('Geographic database management', 'gis'),
  ('SQL', 'database'),
  ('Database management', 'database'),
  ('Data visualization', 'data'),
  ('Data quality', 'data'),
  ('Data analysis', 'data'),
  ('Data science', 'data'),
  ('Machine learning', 'data'),
  ('Deep learning', 'data'),
  ('Data auditing', 'data'),
  ('Reporting', 'data'),
  ('Microsoft Excel', 'office'),
  ('Microsoft Office', 'office'),
  ('Python', 'programming'),
  ('R', 'programming'),
  ('Data protection', 'governance'),
  ('Training', 'soft'),
  ('Documentation', 'soft'),
  ('Technical documentation', 'soft'),
  ('Communication', 'soft'),
  ('Confidentiality', 'soft'),
  ('Networking', 'networking'),
  ('TCP/IP', 'networking'),
  ('VLAN', 'networking'),
  ('LAN/WAN', 'networking'),
  ('Wi-Fi', 'networking'),
  ('Network switching', 'networking'),
  ('Routing', 'networking'),
  ('Router', 'networking'),
  ('Network switch', 'networking'),
  ('CCTV', 'networking'),
  ('Video conferencing', 'networking'),
  ('IT support', 'support'),
  ('Hardware troubleshooting', 'support'),
  ('Troubleshooting', 'support'),
  ('Operating systems', 'support'),
  ('Software installation', 'support'),
  ('IT asset management', 'support'),
  ('Vendor management', 'support'),
  ('Server administration', 'infrastructure'),
  ('Email administration', 'infrastructure'),
  ('Data backup', 'infrastructure'),
  ('Disaster recovery', 'infrastructure'),
  ('Cloud solutions', 'infrastructure'),
  ('Cybersecurity', 'security'),
  ('Access management', 'security'),
  ('Endpoint security', 'security'),
  ('Antivirus', 'security'),
  ('Microsoft Windows', 'platform'),
  ('Microsoft 365', 'platform'),
  ('Microsoft Azure', 'platform'),
  ('Microsoft certification', 'certification');

INSERT OR IGNORE INTO job(
  id, title, employer_id, location, employment_type, role_family, industry, seniority,
  is_ict, education_level, experience_min_years, experience_max_years,
  salary_min, salary_max, salary_currency, salary_period,
  posted_on, deadline, summary, responsibilities, requirements, source_url
) VALUES (
  1,
  'Information Systems & Data Management Officer',
  1,
  'Lilongwe',
  'full_time',
  'data',
  'NGO / humanitarian',
  'mid',
  1,
  'bachelor',
  3,
  NULL,
  NULL, NULL, NULL, NULL,
  '2026-11-13',
  '2026-11-22',
  'Strengthens programme data systems and digital data management, leading DHIS2 configuration, database and GIS support, dashboards and data quality for programme and MEAL teams.',
  '["Design, configure and maintain databases and systems used for programme data.","Configure and maintain DHIS2 metadata, organisation units, data elements, indicators, datasets, validation rules, user roles and dashboards.","Develop and maintain digital data collection forms and workflows using KoboToolbox, CommCare and DHIS2 Mobile.","Develop automated dashboards, visualizations and routine reports using Power BI, Tableau or comparable tools.","Develop and maintain GIS datasets, conduct spatial analysis and produce thematic maps.","Implement data collection standards, validation rules and routine data quality checks.","Train and support programme, MEAL and field staff on data systems, dashboards and GIS outputs.","Provide limited troubleshooting for device, connectivity or software issues affecting programme data systems."]',
  '["Bachelor degree in Information Systems, Data Management, Computer Science, Information Technology, Statistics, GIS, Monitoring and Evaluation or a closely related field.","Specialized training or certification in DHIS2 and GIS.","At least three years of relevant experience in data management, information systems, digital data collection, database administration or programme data systems.","Practical experience configuring or administering DHIS2, including metadata, indicators, validation rules, user roles or dashboards.","Experience with GIS software such as QGIS or ArcGIS, spatial analysis and map production.","Knowledge of database design, data structures, data cleaning and SQL, including SQL Server or MySQL.","Experience developing dashboards and data visualizations using Power BI, Tableau or similar tools.","Experience supporting and training non-technical users.","Experience in an NGO, humanitarian or development programme environment is an advantage.","Python or R is an advantage."]',
  NULL
);

INSERT OR IGNORE INTO job(
  id, title, employer_id, location, employment_type, role_family, industry, seniority,
  is_ict, education_level, experience_min_years, experience_max_years,
  salary_min, salary_max, salary_currency, salary_period,
  posted_on, deadline, summary, responsibilities, requirements, source_url
) VALUES (
  2,
  'GIS Specialist - Agroforestry',
  2,
  'Malawi (International)',
  'full_time',
  'gis',
  'forestry / environmental services',
  'senior',
  1,
  'bachelor',
  3,
  5,
  NULL, NULL, NULL, NULL,
  '2026-11-11',
  '2026-11-18',
  'Transforms complex spatial data into operational products and decision support for a large-scale forest restoration and agroforestry project, combining GIS, remote sensing, programming and data science.',
  '["Manage, structure and assure the quality of geographic databases, including vector, raster, LiDAR and field observation data.","Produce maps, spatial analyses and indicators for planning, monitoring, reporting and decision-making.","Process satellite imagery for land use and land cover classification, forest monitoring, change detection, biomass assessment and environmental evaluation.","Apply machine learning and deep learning techniques to satellite and drone imagery.","Support UAV mission planning, data acquisition and processing of orthomosaics, digital models and point clouds.","Design and maintain digital field data collection tools using Survey123 and ArcGIS ecosystem solutions.","Support forest inventories and MRV or dMRV systems.","Prepare technical reports, presentations, procedures and methodological documentation.","Train and support field teams and non-technical users on GIS outputs."]',
  '["Bachelor or Master degree in Geographic Engineering, Geomatics, Geography, Remote Sensing, Environmental Sciences, Forestry Engineering, Data Science, Computer Science or a related discipline.","Relevant professional experience in GIS, remote sensing and geospatial analysis.","Three to five years of experience is preferred.","Proficiency in ArcGIS Pro, ArcGIS Online and QGIS, including geoprocessing, cartography and geographic database management.","Strong Python skills for data processing, automation and spatial analysis.","Practical experience with Google Earth Engine.","Knowledge of optical and radar imagery, classification methods, spectral indices, time-series analysis and result validation.","Experience in environmental, forestry, conservation, agriculture, carbon or land management projects.","Professional proficiency in English, both written and spoken.","Availability for regular travel to project sites and a valid driving licence."]',
  'https://careers.mota-engil.com/job/Lilongwe-GIS-Specialist-Agroforestry-MF-LI/1456007033/'
);

INSERT OR IGNORE INTO job(
  id, title, employer_id, location, employment_type, role_family, industry, seniority,
  is_ict, education_level, experience_min_years, experience_max_years,
  salary_min, salary_max, salary_currency, salary_period,
  posted_on, deadline, summary, responsibilities, requirements, source_url
) VALUES (
  3,
  'IT Officer',
  3,
  'Blantyre',
  'full_time',
  'infrastructure',
  'logistics / transport',
  'mid',
  1,
  'bachelor',
  2,
  5,
  NULL, NULL, NULL, NULL,
  '2026-11-09',
  '2026-11-20',
  'Provides L1 and L2 IT support and maintains network, server and end-user infrastructure for a multinational logistics operator, including cybersecurity controls and IT asset management.',
  '["Provide effective L1 and L2 technical support to users and ensure timely resolution of IT incidents and service requests.","Support LAN and WAN, Wi-Fi, switching, routing and telecommunications infrastructure.","Administer and support Microsoft-based end-user and identity environments.","Support PCs, laptops, printers, mobile devices and meeting-room or video conferencing equipment.","Monitor infrastructure availability and proactively identify potential service interruptions.","Maintain accurate IT asset configuration records, inventories and technical documentation.","Support cybersecurity controls, access management, endpoint security and compliance with company IT policies.","Liaise with service providers and vendors on incidents, installations and technical projects.","Contribute to continuous improvement, automation and digitalisation initiatives."]',
  '["Bachelor degree in Information Technology, Computer Science, Information Systems, Computer Engineering or a related field.","Preferably two to five years of practical IT support or infrastructure experience.","Good knowledge of Microsoft Windows and Microsoft 365.","Good understanding of networking including TCP/IP, VLANs, LAN and WAN, Wi-Fi, switching and routing.","Practical troubleshooting skills across hardware, software and network environments.","Microsoft certification such as Microsoft 365, Azure or related credentials is an advantage.","Experience supporting CCTV and video conferencing solutions is an advantage.","Exposure to IT infrastructure projects, cybersecurity and endpoint management is an advantage."]',
  NULL
);

INSERT OR IGNORE INTO job(
  id, title, employer_id, location, employment_type, role_family, industry, seniority,
  is_ict, education_level, experience_min_years, experience_max_years,
  salary_min, salary_max, salary_currency, salary_period,
  posted_on, deadline, summary, responsibilities, requirements, source_url
) VALUES (
  4,
  'IT Consultant',
  4,
  'Blantyre',
  'full_time',
  'infrastructure',
  'security services',
  'senior',
  1,
  'diploma',
  3,
  5,
  NULL, NULL, NULL, NULL,
  '2026-10-14',
  NULL,
  'Provides on-demand technical support and maintains hardware, software, internet, server, network and email systems, advising management on IT improvements and technology investments.',
  '["Diagnose and resolve hardware issues affecting computers, laptops, printers, scanners and other IT equipment.","Install, configure, maintain and troubleshoot operating systems, business applications and other software.","Monitor and troubleshoot internet connectivity, Wi-Fi and local area network issues.","Maintain and troubleshoot routers, switches, network cables and other networking equipment.","Maintain, monitor and troubleshoot company servers and related IT infrastructure.","Manage server performance, storage capacity, user accounts, access permissions and shared network resources.","Manage and troubleshoot hosted email services, business email accounts and email configurations.","Ensure regular data backups and verify that backup and recovery procedures work effectively.","Implement IT security measures and assist in investigating and resolving security incidents."]',
  '["Diploma or Bachelor degree in Information Technology, Computer Science, Computer Engineering or a related field.","A minimum of three to five years of relevant experience in IT support, systems administration, network management or a similar role.","Strong knowledge of computer hardware, operating systems, business software, internet connectivity and network troubleshooting.","Practical experience in server administration, hosted email services, data backups and IT infrastructure maintenance.","Good understanding of cybersecurity, data protection, user access management and disaster recovery procedures.","Proficiency in diagnosing and resolving hardware, software, network, server and email-related problems.","Ability to work independently and respond to urgent technical issues.","Willingness to attend the company premises when physical intervention is required."]',
  NULL
);

INSERT OR IGNORE INTO job(
  id, title, employer_id, location, employment_type, role_family, industry, seniority,
  is_ict, education_level, experience_min_years, experience_max_years,
  salary_min, salary_max, salary_currency, salary_period,
  posted_on, deadline, summary, responsibilities, requirements, source_url
) VALUES (
  5,
  'Data Administrator',
  5,
  'Lilongwe',
  'contract',
  'data',
  'FMCG / retail services',
  'mid',
  1,
  'diploma',
  2,
  3,
  NULL, NULL, NULL, NULL,
  '2026-09-21',
  NULL,
  'Maintains the accuracy and integrity of field data and operational databases, performs data audits and analysis, and delivers reporting and insights for operational teams.',
  '["Manage structured reporting schedules and review operational reports for quality.","Perform data analysis to identify exceptions, trends and operational insights.","Conduct time-in-motion studies and consolidate and analyse data as required.","Maintain comprehensive historical data records for auditing and operational tracking.","Manage and maintain operational databases to ensure smooth system functionality.","Conduct data audits, track error patterns and implement systemic improvements.","Collaborate with field teams and software developers to resolve mobility tool or system issues.","Track daily audit activities and calculate field team performance incentives.","Guide team members on policies, processes and data collection standards."]',
  '["Diploma or certification in Information Technology, Data Analytics, Business Administration or a related field.","Two to three years of hands-on experience in a systems or analytical environment.","Experience within the FMCG sector is preferred.","High proficiency in MS Office, especially advanced Excel and data tools.","Strong analytical and critical thinking.","Quality and detail excellence.","Excellent written and verbal communication.","Strong interpersonal and team collaboration skills."]',
  NULL
);

INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Leads configuration and maintenance of DHIS2' FROM skill WHERE name = 'DHIS2';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Develops forms and workflows in KoboToolbox' FROM skill WHERE name = 'KoboToolbox';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Develops digital data collection forms in CommCare' FROM skill WHERE name = 'CommCare';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Develops dashboards using Power BI or Tableau' FROM skill WHERE name = 'Power BI';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Develops dashboards using Power BI or Tableau' FROM skill WHERE name = 'Tableau';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Supports GIS analysis and spatial datasets' FROM skill WHERE name = 'GIS';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'GIS software such as QGIS or ArcGIS' FROM skill WHERE name = 'QGIS';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'GIS software such as QGIS or ArcGIS' FROM skill WHERE name = 'ArcGIS';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Knowledge of SQL and database design' FROM skill WHERE name = 'SQL';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Designs and maintains programme databases' FROM skill WHERE name = 'Database management';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Develops dashboards and data visualizations' FROM skill WHERE name = 'Data visualization';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Implements data quality checks and validation rules' FROM skill WHERE name = 'Data quality';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Analyses programme data for trends and gaps' FROM skill WHERE name = 'Data analysis';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Advanced proficiency in Excel' FROM skill WHERE name = 'Microsoft Excel';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Applies data privacy and access control requirements' FROM skill WHERE name = 'Data protection';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Trains and coaches programme and field staff' FROM skill WHERE name = 'Training';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'required',  'Maintains system configuration documentation' FROM skill WHERE name = 'Documentation';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'preferred', 'Python or R is an advantage' FROM skill WHERE name = 'Python';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 1, id, 'preferred', 'Python or R is an advantage' FROM skill WHERE name = 'R';

INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Strong technical expertise in geographic information systems' FROM skill WHERE name = 'GIS';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Technical expertise in remote sensing' FROM skill WHERE name = 'Remote sensing';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Produces spatial analyses and indicators' FROM skill WHERE name = 'Spatial analysis';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Produces maps and cartographic products' FROM skill WHERE name = 'Cartography';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Proficiency in ArcGIS Pro' FROM skill WHERE name = 'ArcGIS Pro';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Proficiency in ArcGIS Online' FROM skill WHERE name = 'ArcGIS Online';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Proficiency in QGIS' FROM skill WHERE name = 'QGIS';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Strong Python skills for data processing and automation' FROM skill WHERE name = 'Python';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Practical experience with Google Earth Engine' FROM skill WHERE name = 'Google Earth Engine';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Processes satellite imagery for land use classification' FROM skill WHERE name = 'Satellite imagery';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Manages LiDAR and point cloud data' FROM skill WHERE name = 'LiDAR';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Applies machine learning to satellite and drone imagery' FROM skill WHERE name = 'Machine learning';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Applies deep learning where appropriate' FROM skill WHERE name = 'Deep learning';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Designs field data collection tools in Survey123' FROM skill WHERE name = 'Survey123';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Combines GIS with programming and data science' FROM skill WHERE name = 'Data science';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'required',  'Manages and assures quality of geographic databases' FROM skill WHERE name = 'Geographic database management';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'preferred', 'Power BI and GIS integration is preferred' FROM skill WHERE name = 'Power BI';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'preferred', 'UAV mission planning and drone operations' FROM skill WHERE name = 'Drone operations';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'preferred', 'Supports MRV and dMRV systems' FROM skill WHERE name = 'MRV';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'preferred', 'Experience with carbon projects' FROM skill WHERE name = 'Carbon projects';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 2, id, 'preferred', 'Biomass assessment and estimation' FROM skill WHERE name = 'Biomass estimation';

INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Provides L1 and L2 technical support' FROM skill WHERE name = 'IT support';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Good knowledge of Microsoft Windows' FROM skill WHERE name = 'Microsoft Windows';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Good knowledge of Microsoft 365' FROM skill WHERE name = 'Microsoft 365';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Understanding of networking including TCP/IP and VLANs' FROM skill WHERE name = 'Networking';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Good understanding of TCP/IP' FROM skill WHERE name = 'TCP/IP';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Good understanding of VLANs' FROM skill WHERE name = 'VLAN';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Supports LAN and WAN infrastructure' FROM skill WHERE name = 'LAN/WAN';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Supports Wi-Fi infrastructure' FROM skill WHERE name = 'Wi-Fi';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Supports switching infrastructure' FROM skill WHERE name = 'Network switching';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Supports routing infrastructure' FROM skill WHERE name = 'Routing';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Practical troubleshooting of hardware' FROM skill WHERE name = 'Hardware troubleshooting';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Supports cybersecurity controls' FROM skill WHERE name = 'Cybersecurity';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Supports access management' FROM skill WHERE name = 'Access management';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Supports endpoint security' FROM skill WHERE name = 'Endpoint security';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Maintains IT asset configuration records' FROM skill WHERE name = 'IT asset management';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'required',  'Maintains technical documentation' FROM skill WHERE name = 'Technical documentation';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'preferred', 'Azure certification is an added advantage' FROM skill WHERE name = 'Microsoft Azure';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'preferred', 'Microsoft certification is an added advantage' FROM skill WHERE name = 'Microsoft certification';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'preferred', 'Experience supporting CCTV is an advantage' FROM skill WHERE name = 'CCTV';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'preferred', 'Experience supporting video conferencing is an advantage' FROM skill WHERE name = 'Video conferencing';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 3, id, 'preferred', 'Liaises with service providers and vendors' FROM skill WHERE name = 'Vendor management';

INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Provides IT support on an as-needed basis' FROM skill WHERE name = 'IT support';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Diagnoses and resolves hardware issues' FROM skill WHERE name = 'Hardware troubleshooting';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Installs and maintains operating systems' FROM skill WHERE name = 'Operating systems';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Installs and configures business applications' FROM skill WHERE name = 'Software installation';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Maintains network and internet connectivity' FROM skill WHERE name = 'Networking';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Troubleshoots local area network issues' FROM skill WHERE name = 'LAN/WAN';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Troubleshoots Wi-Fi connectivity' FROM skill WHERE name = 'Wi-Fi';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Maintains and troubleshoots routers' FROM skill WHERE name = 'Router';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Maintains and troubleshoots switches' FROM skill WHERE name = 'Network switch';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Maintains and monitors company servers' FROM skill WHERE name = 'Server administration';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Manages hosted email services' FROM skill WHERE name = 'Email administration';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Ensures regular data backups' FROM skill WHERE name = 'Data backup';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Ensures appropriate security measures' FROM skill WHERE name = 'Cybersecurity';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Installs and maintains antivirus software' FROM skill WHERE name = 'Antivirus';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Manages user accounts and access permissions' FROM skill WHERE name = 'Access management';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Maintains backup and recovery procedures' FROM skill WHERE name = 'Disaster recovery';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Diagnoses and resolves technical issues' FROM skill WHERE name = 'Troubleshooting';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'High levels of confidentiality' FROM skill WHERE name = 'Confidentiality';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'required',  'Good understanding of data protection' FROM skill WHERE name = 'Data protection';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'preferred', 'Recommends cloud-based solutions' FROM skill WHERE name = 'Cloud solutions';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 4, id, 'preferred', 'Liaises with internet service providers' FROM skill WHERE name = 'Vendor management';

INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'required',  'Advanced Excel and data tools' FROM skill WHERE name = 'Microsoft Excel';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'required',  'High proficiency in MS Office' FROM skill WHERE name = 'Microsoft Office';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'required',  'Performs data analysis to identify trends' FROM skill WHERE name = 'Data analysis';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'required',  'Ensures credibility, accuracy and integrity of field data' FROM skill WHERE name = 'Data quality';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'required',  'Manages and maintains operational databases' FROM skill WHERE name = 'Database management';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'required',  'Manages structured reporting schedules' FROM skill WHERE name = 'Reporting';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'required',  'Conducts data audits and tracks error patterns' FROM skill WHERE name = 'Data auditing';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'required',  'Excellent written and verbal communication' FROM skill WHERE name = 'Communication';
INSERT OR IGNORE INTO job_skill(job_id, skill_id, requirement, evidence) SELECT 5, id, 'mentioned', 'Guides team members on policies and processes' FROM skill WHERE name = 'Training';
