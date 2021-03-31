# go-skydisc
Optimization of operations over resources and time

The following elements are managed

Service area
describes e.g. a branch office, temporal allocation of certain resources
* Designation
* Address, geo-coordinates for starting point
* Capacities per qualification

Order
Order with the following properties
* class of service (e.g. unplanned ticket, planned maintenance)
* describes a service
* priority, distress
* Project number and name
* earliest start and latest end - e.g. agreed deadline
* predecessor relationships
* Services location: Address, contact person, client, geo-coordinates.
* requirement: necessary qualification, trade and number of resources
* Duration
* Service area (e.g. branch office)

Resources
Employees, machines or equipment
* Name, designation, telephone, email
* Qualification
* Working time, absence, calendar
* Service area (e.g. branch)
* Substitution

Assignment/Appointment
a resource can be booked for an order with one or more appointments
* Plan: start, end, duration
* Actual: Start, end, duration
* set up time, arrival time
* Probability of realization
* necessary predecessors

Catalogs:
* Trade and qualifications 
* 
 
