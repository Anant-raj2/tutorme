class Person {
  constructor(name, gender, age) {
    this.name = name;
    this.gender = gender;
    this.age = age;
  }
}

class Tutor extends Person {
  constructor(name, gender, age, subjects, availability) {
    super(name, gender, age);
    this.subjects = subjects;
    this.availability = availability; // hours per week
  }
}

class Student extends Person {
  constructor(name, gender, age, wantedSubjects, preferredGender, availableHours) {
    super(name, gender, age);
    this.wantedSubjects = wantedSubjects;
    this.preferredGender = preferredGender;
    this.availableHours = availableHours;
  }
}

class TutoringMatch {
  constructor(weights = {
    subjectMatch: 0.4,
    genderPreference: 0.2,
    ageProximity: 0.2,
    availability: 0.2
  }) {
    this.weights = weights;
  }

  findBestMatch(student, tutors) {
    let bestMatch = null;
    let bestScore = -Infinity;

    for (const tutor of tutors) {
      const score = this.calculateMatchScore(student, tutor);
      if (score > bestScore) {
        bestScore = score;
        bestMatch = tutor;
      }
    }

    return { tutor: bestMatch, score: bestScore };
  }

  calculateMatchScore(student, tutor) {
    const subjectMatchScore = this.calculateSubjectMatchScore(student, tutor);
    const genderPreferenceScore = this.calculateGenderPreferenceScore(student, tutor);
    const ageProximityScore = this.calculateAgeProximityScore(student, tutor);
    const availabilityScore = this.calculateAvailabilityScore(student, tutor);

    return (
      subjectMatchScore * this.weights.subjectMatch +
      genderPreferenceScore * this.weights.genderPreference +
      ageProximityScore * this.weights.ageProximity +
      availabilityScore * this.weights.availability
    );
  }

  calculateSubjectMatchScore(student, tutor) {
    const matchedSubjects = student.wantedSubjects.filter(subject =>
      tutor.subjects.includes(subject)
    );
    return matchedSubjects.length / student.wantedSubjects.length;
  }

  calculateGenderPreferenceScore(student, tutor) {
    return student.preferredGender === tutor.gender ? 1 : 0;
  }

  calculateAgeProximityScore(student, tutor) {
    const ageDifference = Math.abs(student.age - tutor.age);
    return Math.max(0, 1 - ageDifference / 50); // Assume max age difference of 50 years
  }

  calculateAvailabilityScore(student, tutor) {
    const minHours = Math.min(student.availableHours, tutor.availability);
    return minHours / student.availableHours;
  }
}

// Example usage
const tutors = [
  new Tutor("Alice", "female", 28, ["Math", "Physics", "Chemistry"], 20),
  new Tutor("Bob", "male", 35, ["English", "History", "Literature"], 15),
  new Tutor("Charlie", "male", 42, ["Computer Science", "Math"], 25),
  new Tutor("Diana", "female", 31, ["Biology", "Chemistry"], 30),
];

const students = [
  new Student("Eve", "female", 16, ["Math", "Physics"], "female", 10),
  new Student("Frank", "male", 18, ["English", "History"], null, 5),
  new Student("Grace", "female", 17, ["Computer Science", "Math"], "male", 15),
];

const matcher = new TutoringMatch();

for (const student of students) {
  const { tutor, score } = matcher.findBestMatch(student, tutors);
  console.log(`Best match for ${student.name}: ${tutor.name} (Score: ${score.toFixed(2)})`);
}
