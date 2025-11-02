declare namespace API {
  type BaseResponse = {
    code: string;
    message: string;
  };

  type getHostParams = {
    domain?: string;
  };

  type getOptChangeParams = {
    type?: string;
  };

  type getOptParams = {
    type?: string;
  };

  type HostVo = {
    domain: string;
    ip: string;
    type: string;
  };

  type OptRequest = {
    type: string;
    data: OptVo[];
  };

  type OptVo = {
    ip: string;
    delay: string;
    rate: string;
  };

  type QueryHostListResponse = {
    code: number;
    message: string;
    data: { total: number; list: HostVo[] };
  };

  type QueryHostResponse = {
    code: number;
    message: string;
    data: HostVo;
  };
}
