export type Requirement = 'required' | 'preferred' | 'mentioned'

export type JobSkill = {
  name: string
  category: string
  requirement: Requirement
  evidence: string
}

export type Range = {
  min: number | null
  max: number | null
}

export type Salary = {
  disclosed: boolean
  min: number | null
  max: number | null
  currency: string | null
  period: string | null
}

export type Job = {
  id: number
  title: string
  employer: string
  location: string
  employment_type: string
  role_family: string
  industry: string | null
  seniority: string
  is_ict: boolean
  education_level: string | null
  experience_years: Range
  salary: Salary
  posted_on: string
  deadline: string | null
  summary: string
  responsibilities: string[]
  requirements: string[]
  skills: JobSkill[]
}

export type SkillDemand = {
  name: string
  category: string
  jobs: number
  required: number
}

export type JobQuery = {
  roleFamily?: string
  skill?: string
  postedFrom?: string
  postedTo?: string
  limit?: number
  offset?: number
}

export type JobListResponse = {
  jobs: Job[]
  count: number
  limit: number
  offset: number
}

export type SkillListResponse = {
  skills: SkillDemand[]
  count: number
}

const baseURL = import.meta.env.VITE_API_BASE ?? ''

async function request<T>(path: string, signal?: AbortSignal): Promise<T> {
  const response = await fetch(`${baseURL}${path}`, {
    signal,
    headers: { Accept: 'application/json' },
  })

  if (!response.ok) {
    let message = `Request failed (${response.status})`
    try {
      const body = (await response.json()) as { error?: string }
      if (body.error) {
        message = body.error
      }
    } catch {
      // response had no JSON body
    }
    throw new Error(message)
  }

  return (await response.json()) as T
}

export function listJobs(query: JobQuery, signal?: AbortSignal): Promise<JobListResponse> {
  const params = new URLSearchParams()
  if (query.roleFamily) params.set('role_family', query.roleFamily)
  if (query.skill) params.set('skill', query.skill)
  if (query.postedFrom) params.set('posted_from', query.postedFrom)
  if (query.postedTo) params.set('posted_to', query.postedTo)
  params.set('limit', String(query.limit ?? 24))
  params.set('offset', String(query.offset ?? 0))

  return request<JobListResponse>(`/api/jobs?${params.toString()}`, signal)
}

export function listSkills(signal?: AbortSignal): Promise<SkillListResponse> {
  return request<SkillListResponse>('/api/skills', signal)
}

export const roleFamilies = [
  { value: '', label: 'All role families' },
  { value: 'software', label: 'Software' },
  { value: 'data', label: 'Data' },
  { value: 'gis', label: 'GIS' },
  { value: 'infrastructure', label: 'Infrastructure' },
  { value: 'support', label: 'Support' },
  { value: 'security', label: 'Security' },
  { value: 'management', label: 'Management' },
  { value: 'other', label: 'Other' },
] as const
